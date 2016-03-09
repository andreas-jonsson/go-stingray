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
	"io"
	"math"
	"math/rand"
	"os"
	"time"

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

	rand.Seed(time.Now().UnixNano())
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

func transferThumbnail(con *console.Console, writer io.Writer, id int) {
	defer func() {
		if r := recover(); r != nil {
			errorln("received corrupt data")
		}
	}()

	for {
		_, data, err := con.Receive()
		assertln(err, err)

		if len(data) == 0 {
			continue
		}

		lex := sjson.NewLexer(bytes.NewBuffer(data))
		obj, err := sjson.Decode(lex)
		assertln(err, err)

		m := obj.(map[string]sjson.Value)
		if m["type"].(string) != "thumbnail" || int(m["id"].(float64)) != id {
			continue
		}

		_, err = lex.Reader().WriteTo(writer)
		assertln(err, err)
		return
	}
}

func transferJittered(con *console.Console) *frameCapture {
	defer func() {
		if r := recover(); r != nil {
			errorln("received corrupt data")
		}
	}()

	var capture *frameCapture
	for capture == nil || !capture.isComplete() {
		_, data, err := con.Receive()
		assertln(err, err)

		if len(data) == 0 {
			continue
		}

		lex := sjson.NewLexer(bytes.NewBuffer(data))
		obj, err := sjson.Decode(lex)
		assertln(err, err)

		m := obj.(map[string]sjson.Value)
		if m["type"].(string) != "frame_capture" {
			continue
		}

		reader := lex.Reader()
		b, err := reader.ReadByte()
		assertln(err, err)
		if b != 0 {
			errorln("invalid binary message")
		}

		id := int(m["id"].(float64))
		tap := int(m["tap"].(float64))
		format, _ := m["surface_format"].(string)

		if format != "R8G8B8A8" {
			errorln("invalid surface format: " + format)
		}

		if capture == nil {
			fmt.Println("transfer image...")

			numTaps := int(m["num_taps"].(float64))
			stride := int(m["stride"].(float64))
			capture = newFrameCapture(id, numTaps, stride)
		}

		var buf bytes.Buffer
		_, err = reader.WriteTo(&buf)
		assertln(err, err)

		if err := capture.addTap(tap, buf.Bytes()); err != nil {
			errorln(err)
		}
	}

	return capture
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

	if arguments.scale == 1 {
		id := rand.Intn(math.MaxInt32)
		cmd := fmt.Sprintf("FrameCapture.thumbnail(ConsoleServer.current_client_id(),nil,'back_buffer',%v,Renderer.back_buffer_size())", id)
		err = con.SendCommand(console.Script, cmd)
		assertln(err, err)

		fp, err := os.Create(arguments.outputPath)
		assertln(err, err)
		defer fp.Close()

		fmt.Println("waiting for thumbnail...")
		transferThumbnail(con, fp, id)
	} else {
		cmd := fmt.Sprintf("FrameCapture.replay_jittered_frame('console_send',nil,%v,nil)", arguments.scale)
		err = con.SendCommand(console.Script, cmd)
		assertln(err, err)

		fmt.Println("waiting for taps...")
		capture := transferJittered(con)

		if err := capture.save(arguments.outputPath); err != nil {
			errorln(err)
		}
	}

	fmt.Println("complete")
	fmt.Println("output written to: " + arguments.outputPath)
}
