/*
Copyright (C) 2016 Andreas T Jonsson

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
	"bytes"
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/andreas-jonsson/go-stingray/console"
	"github.com/andreas-jonsson/go-stingray/sjson"
)

var arguments struct {
	hostAddress,
	outputPath string
	scale int
}

func init() {
	flag.Usage = func() {
		fmt.Printf("Usage: screenshot [options]\n\n")
		flag.PrintDefaults()
	}

	flag.StringVar(&arguments.hostAddress, "host", "localhost", "host address, address:[port]")
	flag.StringVar(&arguments.outputPath, "output", "screenshot.png", "write image to file")
	flag.IntVar(&arguments.scale, "scale", 1, "screen-buffer multiplier")
}

func errorln(msg ...interface{}) {
	fmt.Fprintln(os.Stderr, msg...)
	os.Exit(-1)
}

func assertln(err error, msg ...interface{}) {
	if err != nil {
		errorln(msg)
	}
}

func main() {
	flag.Parse()
	fmt.Println("Stingray Hi-Res Screenshot")
	fmt.Printf("Copyright (C) 2016 Andreas T Jonsson\n\n")

	if arguments.scale < 1 || arguments.scale > 32 {
		errorln("invalid scale value")
	}

	fmt.Printf("connecting to %s...\n", arguments.hostAddress)
	con, err := console.NewConsole(arguments.hostAddress, "")
	assertln(err, "could not connect to: "+arguments.hostAddress)
	defer con.Close()

	fmt.Println("connected")

	cmd := fmt.Sprintf("FrameCapture.replay_jittered_frame('console_send',nil,%v,nil)", arguments.scale)
	fmt.Println(cmd)
	sjson.Encode(os.Stdout, cmd)

	err = con.SendCommand(console.Script, cmd)
	assertln(err, err)

	defer func() {
		if r := recover(); r != nil {
			errorln("received corrupt data")
		}
	}()

	var capture *frameCapture
	fmt.Println("waiting for response...")

	for {
		_, data, err := con.Receive()
		assertln(err, err)

		if len(data) == 0 {
			continue
		}

		reader := bytes.NewBuffer(data)
		obj, err := sjson.Decode(sjson.NewLexer(reader))
		assertln(err, err)

		m := obj.(map[string]sjson.Value)
		if m["type"].(string) != "frame_capture" {
			continue
		}

		b, err := reader.ReadByte()
		assertln(err, err)
		if b != 0 {
			errorln("invalid binary message")
		}

		id, _ := strconv.Atoi(m["id"].(string))
		tap, _ := strconv.Atoi(m["tap"].(string))
		format, _ := m["surface_format"].(string)

		if format != "R8G8B8A8" {
			errorln("invalid surface format: " + format)
		}

		if capture == nil {
			fmt.Println("transfer image...")

			numTaps, _ := strconv.Atoi(m["num_taps"].(string))
			stride, _ := strconv.Atoi(m["stride"].(string))
			capture = newFrameCapture(id, numTaps, stride)
		}

		if err := capture.addTap(tap, reader.Bytes()); err != nil {
			errorln(err)
		}

		if capture.isComplete() {
			if err := capture.save(arguments.outputPath); err != nil {
				errorln(err)
			}

			fmt.Println("compleat")
			return
		}
	}
}
