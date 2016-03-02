/*
Copyright (C) 2015-2016 Andreas T Jonsson

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/andreas-jonsson/go-stingray/data-server"
	"github.com/andreas-jonsson/go-stingray/data-server/storage"
)

const (
	versionString     = "2.1.0"
	defaultBufferSize = 10 * 1024 * 1034
)

var (
	database      storage.Database
	running       int32 = 1
	driverFlag    string
	sourceFlag    string
	locklessFlag  bool
	broadcastFlag bool
)

func init() {
	flag.Usage = func() {
		fmt.Printf("Usage: server [options]\n\n")
		flag.PrintDefaults()
	}

	flag.StringVar(&driverFlag, "driver", "sqlite3", "driver type, (sqlite3, mysql, bolt, ram)")
	flag.StringVar(&sourceFlag, "source", "cache.db", "database source specifier")
	flag.BoolVar(&locklessFlag, "lockless", false, "access database from multiple threads")
	flag.BoolVar(&broadcastFlag, "broadcast", true, "respond to broadcast messages")
}

func main() {
	fmt.Printf("Stingray Cache Server v %s\n", versionString)
	fmt.Printf("Copyright (C) 2015-2016 Andreas T Jonsson\n\n")

	flag.Parse()
	setupSignals()

	log.Println("opening database:", driverFlag, sourceFlag, locklessFlag)
	if driverFlag == "ram" {
		database = storage.NewRAMDatabase()
	} else if driverFlag == "bolt" {
		database = storage.NewBoltDatabase(sourceFlag)
	} else {
		database = storage.NewSQLDatabase(driverFlag, sourceFlag, locklessFlag)
	}

	defer func() {
		database.Close()
		log.Println("database closed")
	}()

	udpCon, tcpList := setupListeners(server.BroadcastPort, server.DefaultPort)
	defer tcpList.Close()
	defer udpCon.Close()

	if broadcastFlag {
		go listenForBroadcasts(udpCon)
	}

	listenForConnections(tcpList)
}

func clientUpload(s storage.Session, tcpCon *net.TCPConn, buffer []byte, requestHeader [2]uint64) ([]byte, error) {
	hash := requestHeader[0]
	size := requestHeader[1]

	veryBig := uint64(250 * 1024 * 1024)
	if size > veryBig {
		return buffer, fmt.Errorf("abnormal file size, %d MB", size/1024/1024)
	}

	for oldCap := cap(buffer); oldCap < int(size); {
		oldCap = oldCap*2 + 1
		buffer = make([]byte, 0, oldCap)
	}

	inputBuffer := buffer[:size]

	for written := 0; written < int(size); {
		tcpCon.SetDeadline(time.Now().Add(time.Second))
		sz, err := tcpCon.Read(inputBuffer[written:size])
		if err != nil {
			e, ok := err.(net.Error)
			if !ok || !e.Timeout() {
				return buffer, err
			}
		}
		written += sz
	}

	s.Store(hash, inputBuffer)
	return buffer, nil
}

func clientDownload(s storage.Session, tcpCon *net.TCPConn, buffer []byte, requestHeader [2]uint64) ([]byte, error) {
	hash := requestHeader[0]
	buffer = s.Load(hash, buffer)

	size := len(buffer)
	requestHeader[1] = uint64(size)

	tcpCon.SetDeadline(time.Now().Add(time.Second))
	err := binary.Write(tcpCon, binary.LittleEndian, &requestHeader)
	if err != nil {
		return buffer, err
	}

	for written := 0; written < size; {
		tcpCon.SetDeadline(time.Now().Add(time.Second))
		sz, err := tcpCon.Write(buffer[written:])
		if err != nil {
			e, ok := err.(net.Error)
			if !ok || !e.Timeout() {
				return buffer, err
			}
		}
		written += sz
	}

	return buffer, nil
}

func serveClient(tcpCon *net.TCPConn) {
	defer tcpCon.Close()

	headerSize := len(server.ProtocolHeader) + 1
	headerBuffer := make([]byte, headerSize)

	tcpCon.SetDeadline(time.Now().Add(time.Second))
	size, err := tcpCon.Read(headerBuffer)

	ip := tcpCon.RemoteAddr().String()
	if err != nil || string(headerBuffer[:size-1]) != server.ProtocolHeader {
		log.Println("could not validate protocol used by", ip)
	} else {
		log.Println("connection established to", ip)
		buffer := make([]byte, 0, defaultBufferSize)

		log.Println("open database session for", ip)
		session, err := database.Open()

		if err == nil {
			defer session.Close()

			for atomic.LoadInt32(&running) == 1 {
				var requestHeader [2]uint64

				tcpCon.SetDeadline(time.Now().Add(time.Minute * 15))
				if err := binary.Read(tcpCon, binary.LittleEndian, &requestHeader); err != nil {
					e, ok := err.(net.Error)
					if ok && e.Timeout() {
						continue
					}
					break
				}

				if requestHeader[1] > 0 {
					buffer, err = clientUpload(session, tcpCon, buffer, requestHeader)
				} else {
					buffer, err = clientDownload(session, tcpCon, buffer, requestHeader)
				}

				if err != nil {
					log.Println(err)
					break
				}
			}
		}

		log.Println("closing connection to", ip)
	}
}

func listenForConnections(tcpList *net.TCPListener) {
	for atomic.LoadInt32(&running) == 1 {
		tcpList.SetDeadline(time.Now().Add(time.Second))
		tcpCon, err := tcpList.AcceptTCP()
		if err == nil {
			go serveClient(tcpCon)
		}
	}
}

func listenForBroadcasts(udpCon *net.UDPConn) {
	headerSize := len(server.ProtocolHeader) + 1
	buffer := make([]byte, headerSize)
	for atomic.LoadInt32(&running) == 1 {
		udpCon.SetDeadline(time.Now().Add(time.Second))
		if size, remote, err := udpCon.ReadFromUDP(buffer[:]); err == nil {
			if magic := string(buffer[:size-1]); magic == server.ProtocolHeader {
				go sendBroadcastResponse(remote)
			} else {
				log.Println(remote, "invalid protocol", magic)
			}
		}
	}
}

func sendBroadcastResponse(addr *net.UDPAddr) {
	log.Printf("reply on broadcast to %s:%d\n", addr.IP, server.InfoPort)
	if udpCon, err := net.Dial("udp", fmt.Sprintf("%s:%d", addr.IP, server.InfoPort)); err != nil {
		log.Println(err)
	} else {
		data := []byte(server.ProtocolHeader)
		data = append(data, 0)
		udpCon.SetDeadline(time.Now().Add(time.Second))
		udpCon.Write(data)
		udpCon.Close()
	}
}

func setupListeners(udpPort int, tcpPort int) (*net.UDPConn, *net.TCPListener) {
	log.Println("listening for broadcast on port", udpPort)
	udpAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", udpPort))
	checkErr(err)

	udpCon, err := net.ListenUDP("udp", udpAddr)
	checkErr(err)

	log.Println("listening for connections on port", tcpPort)
	tcpAddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf(":%d", tcpPort))
	checkErr(err)

	tcpList, err := net.ListenTCP("tcp", tcpAddr)
	checkErr(err)

	return udpCon, tcpList
}

func setupSignals() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	go func() {
		for sig := range c {
			log.Println("recived signal", sig)
			log.Println("closing connections")
			atomic.StoreInt32(&running, 0)
		}
	}()
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
