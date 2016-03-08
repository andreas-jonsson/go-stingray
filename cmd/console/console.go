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
	"flag"
	"fmt"
	"os"

	"github.com/andreas-jonsson/go-stingray/console"
)

var arguments struct {
	hostAddress,
	inputFile string
	quiet bool
}

func init() {
	flag.Usage = func() {
		fmt.Printf("Usage: console [options]\n\n")
		flag.PrintDefaults()
	}

	flag.StringVar(&arguments.hostAddress, "host", "localhost", "host address, address:[port]")
	flag.StringVar(&arguments.inputFile, "input", "", "input file, '-' for stdin")
	flag.BoolVar(&arguments.quiet, "q", false, "quiet, don't print any extra information")
}

func main() {
	flag.Parse()
	q := arguments.quiet
	host := arguments.hostAddress

	if !q {
		fmt.Println("Stingray Console")
		fmt.Printf("Copyright (C) 2016 Andreas T Jonsson\n\n")
		fmt.Printf("connecting to %s...\n", host)
	}

	con, err := console.NewConsole(host, "")
	if err != nil {
		fmt.Fprintln(os.Stderr, "could not connect to: "+host)
		os.Exit(-1)
	}

	defer con.Close()
	if !q {
		fmt.Println("connected")
	}

	if arguments.inputFile != "" {
		go processInput(con)
	}

	for {
		msg, err := con.ReceiveMessage()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(-1)
		}
		fmt.Println(msg)
	}
}

func processInput(con *console.Console) {
	var err error
	fp := os.Stdin

	file := arguments.inputFile
	if file != "-" {
		if fp, err = os.Open(file); err != nil {
			fmt.Fprintln(os.Stderr, "could not open: "+file)
			os.Exit(-1)
		}
	}

	scanner := bufio.NewScanner(fp)
	for scanner.Scan() {
		if line := scanner.Text(); len(line) > 0 {
			ty := console.Command
			if line[0] == '#' {
				ty = console.Script
				line = line[1:]
			}

			if err := con.SendCommand(ty, line); err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(-1)
			}
		}
	}
}
