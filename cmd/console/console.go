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
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/andreas-jonsson/go-stingray/console"
	"github.com/jroimartin/gocui"
)

var arguments struct {
	hostAddress,
	inputFile string
	quiet bool
}

func errorln(msg ...interface{}) {
	if arguments.quiet {
		fmt.Fprintln(os.Stderr, msg...)
		os.Exit(-1)
	} else {
		goCUI.Execute(func(g *gocui.Gui) error {
			g.Close()
			fmt.Fprintln(os.Stderr, msg...)
			os.Exit(-1)
			return nil
		})
		time.Sleep(time.Minute)
	}
}

func assertln(err error, msg ...interface{}) {
	if err != nil {
		errorln(msg...)
	}
}

func assertErrln(err error) {
	if err != nil {
		errorln(err)
	}
}

func println(msg ...interface{}) {
	doPrint("", msg)
}

func printf(f string, msg ...interface{}) {
	doPrint(f, msg)
}

func showLicense() {
	url := "https://raw.githubusercontent.com/andreas-jonsson/go-stingray/master/LICENSE"
	res, err := http.Get(url)
	if err != nil {
		println(url)
		return
	}
	defer res.Body.Close()

	license, err := ioutil.ReadAll(res.Body)
	if err != nil {
		println(url)
		return
	}

	println(string(license))
}

func quiet() {
	host := arguments.hostAddress
	con, err := console.NewConsole(host, "")
	assertln(err, errors.New("could not connect to: "+host))
	defer con.Close()

	if arguments.inputFile != "" {
		go processInput(con)
	}

	for {
		msg, err := con.ReceiveMessage()
		assertErrln(err)
		fmt.Println(msg)
	}
}

func executeCommand(con *console.Console, cmd string) error {
	ty := console.Command
	if cmd[0] == '#' {
		ty = console.Script
		cmd = cmd[1:]
	}
	return con.SendCommand(ty, cmd)
}

func processInput(con *console.Console) {
	var err error
	fp := os.Stdin

	file := arguments.inputFile
	if file != "-" {
		fp, err = os.Open(file)
		assertln(err, errors.New("could not open: "+file))
	}

	scanner := bufio.NewScanner(fp)
	for scanner.Scan() {
		if line := scanner.Text(); len(line) > 0 {
			assertErrln(executeCommand(con, line))
		}
	}
}

func init() {
	flag.Usage = func() {
		fmt.Printf("Usage: console [options]\n\n")
		flag.PrintDefaults()
	}

	flag.StringVar(&arguments.hostAddress, "host", "localhost", "host address, address:[port]")
	flag.StringVar(&arguments.inputFile, "input", "", "input file, '-' for stdin")
	flag.BoolVar(&arguments.quiet, "q", false, "no GUI, pipe-only")
}

func main() {
	flag.Parse()
	if arguments.quiet {
		quiet()
	} else {
		gui()
	}
}
