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
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/andreas-jonsson/go-stingray/console"
	"github.com/jroimartin/gocui"
)

var goCUI *gocui.Gui

func doPrint(f string, msg []interface{}) {
	if !arguments.quiet {
		goCUI.Execute(func(g *gocui.Gui) error {
			v, err := g.View("top")
			if err != nil {
				return err
			}
			fmt.Fprint(v, " ")
			if f == "" {
				fmt.Fprintln(v, msg...)
			} else {
				fmt.Fprintf(v, f, msg...)
			}
			return nil
		})
	}
}

func setTitle(f string, a ...interface{}) {
	goCUI.Execute(func(g *gocui.Gui) error {
		v, _ := g.View("top")
		v.Title = fmt.Sprintf(f, a...)
		return nil
	})
}

func layout(g *gocui.Gui) error {
	const inputHeight = 3

	maxX, maxY := g.Size()
	if v, err := g.SetView("top", 0, 0, maxX-1, maxY-inputHeight); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		fmt.Fprint(v, logo)
		fmt.Fprintf(v, "        Copyright (C) 2016 Andreas T Jonsson\n\n")
		fmt.Fprint(v, notice)

		v.FgColor = gocui.ColorWhite
		v.Title = "disconnected"
		v.Autoscroll = true
		v.Wrap = false
	}

	if v, err := g.SetView("bottom", 0, maxY-inputHeight, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		if err := g.SetCurrentView("bottom"); err != nil {
			return err
		}

		v.FgColor = gocui.ColorWhite
	}

	return nil
}

func setupInputKeybindings(con *console.Console) {
	goCUI.Execute(func(g *gocui.Gui) error {
		enter := func(g *gocui.Gui, v *gocui.View) error {
			str := strings.TrimSpace(v.Buffer())

			switch str {
			case "show c", "show w":
				showLicense()
			default:
				if len(str) > 0 {
					printf("> %s\n", str)
					assertErrln(executeCommand(con, str))
				}
			}

			v.Clear()
			return nil
		}

		v, err := g.View("bottom")
		assertErrln(err)
		v.Editable = true

		assertErrln(g.SetKeybinding("bottom", gocui.KeyEnter, gocui.ModNone, enter))
		return nil
	})
}

func setupKeybindings() {
	quit := func(*gocui.Gui, *gocui.View) error {
		return gocui.ErrQuit
	}
	assertErrln(goCUI.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit))
}

func gui() {
	goCUI = gocui.NewGui()
	if err := goCUI.Init(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(-1)
	}
	defer goCUI.Close()

	goCUI.FgColor = gocui.ColorCyan
	goCUI.BgColor = gocui.ColorBlue
	goCUI.Cursor = true
	goCUI.SetLayout(layout)
	setupKeybindings()

	go func() {
		host := arguments.hostAddress
		printf("connecting to %s...\n", host)

		con, err := console.NewConsole(host, "")
		assertln(err, errors.New("could not connect to host"))
		defer con.Close()

		println("connected")
		setTitle(host)
		setupInputKeybindings(con)

		for {
			msg, err := con.ReceiveMessage()
			assertErrln(err)
			println(msg)

		}
	}()

	if err := goCUI.MainLoop(); err != nil && err != gocui.ErrQuit {
		errorln(err)
	}
}
