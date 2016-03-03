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

var (
	hostAddress,
	inputFile string
)

func init() {
	flag.Usage = func() {
		fmt.Printf("Usage: console [options]\n\n")
		flag.PrintDefaults()
	}

	flag.StringVar(&hostAddress, "host", "localhost", "host address, address:[port]")
	flag.StringVar(&inputFile, "input", "", "input file, '-' for stdin")
}

func main() {
	flag.Parse()
	fmt.Println("Stingray Console")
	fmt.Printf("Copyright (C) 2016 Andreas T Jonsson\n\n")

	fmt.Printf("connecting to %s...\n", hostAddress)
	con, err := console.NewConsole(hostAddress)
	if err != nil {
		fmt.Println("could not connect to: " + hostAddress)
		os.Exit(-1)
	}

	defer con.Close()
	fmt.Println("connected")

	if inputFile != "" {
		go processInput(con)
	}

	for {
		msg, err := con.Read()
		if err != nil {
			con.Close()
			os.Exit(-1)
		}
		fmt.Println(msg)
	}
}

func processInput(con *console.Console) {
	var err error
	fp := os.Stdin

	if inputFile != "-" {
		if fp, err = os.Open(inputFile); err != nil {
			fmt.Println("could not open: " + inputFile)
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

			if err := con.Write(ty, line); err != nil {
				fmt.Println(err)
				os.Exit(-1)
			}
		}
	}
}
