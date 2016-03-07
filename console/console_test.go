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

package console

import (
	"fmt"
	"io"
	"net/http"
	"reflect"
	"testing"

	"github.com/andreas-jonsson/go-stingray/sjson"

	"golang.org/x/net/websocket"
)

func consoleServer(ws *websocket.Conn) {
	io.Copy(ws, ws)
}

func startServer(t *testing.T) {
	http.Handle("/", websocket.Handler(consoleServer))
	if err := http.ListenAndServe(fmt.Sprintf(":%v", DefaultPort), nil); err != nil {
		t.Fatal(err)
	}
}

func receiveAndTest(t *testing.T, con *Console, expected sjson.Value) {
	msg, err := con.Receive()
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(msg, expected) {
		t.Fail()
	}
}

func TestConsole(t *testing.T) {
	go startServer(t)
	con, err := NewConsole("localhost", "")
	if err != nil {
		t.Fatal(err)
	}
	defer con.Close()

	if err := con.SendCommand(Command, "test arg1 arg2 arg3"); err != nil {
		t.Fatal(err)
	}

	args := []sjson.Value{"arg1", "arg2", "arg3"}
	cmd := map[string]sjson.Value{"type": "command", "command": "test", "arg": args}
	receiveAndTest(t, con, cmd)

	if err := con.SendCommand(Script, "test"); err != nil {
		t.Fatal(err)
	}

	cmd = map[string]sjson.Value{"type": "script", "script": "test"}
	receiveAndTest(t, con, cmd)
}
