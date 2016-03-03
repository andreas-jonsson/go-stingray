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
	"bytes"
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/andreas-jonsson/go-stingray/sjson"
	"golang.org/x/net/websocket"
)

type CommandType int

const (
	Command CommandType = iota
	Script
)

const (
	DefaultPort        = 14030
	DefaultXboxOnePort = 4601
)

const (
	socketTimeout = 5
)

type Message struct {
	System,
	Level,
	MessageType,
	Message string
}

func (m Message) String() string {
	if m.System == "" {
		return m.Message
	} else {
		return fmt.Sprintf("[%s] %s", m.System, m.Message)
	}
}

type Console struct {
	lex  *sjson.Lexer
	ws   *websocket.Conn
	host string
}

func (con *Console) Read() (Message, error) {
	var (
		err error
		msg Message
	)

	defer func() {
		if r := recover(); r != nil {
			err = errors.New("invalid message")
		}
	}()

	for {
		val, err := sjson.Decode(con.lex)
		if err != nil {
			return msg, errors.New("could not decode message")
		}

		m := val.(map[string]sjson.Value)
		if m["type"].(string) != "message" {
			continue
		}

		msg.System = m["system"].(string)
		msg.Message = m["message"].(string)
		msg.MessageType = m["message_type"].(string)
		msg.Level = m["level"].(string)
		break
	}

	return msg, err
}

func (con *Console) Write(ty CommandType, command string) error {
	var buf bytes.Buffer

	switch ty {
	case Command:
		var argsValue []sjson.Value

		args := strings.Split(command, " ")
		for _, arg := range args[1:] {
			argsValue = append(argsValue, arg)
		}

		m := map[string]sjson.Value{"type": "command", "command": args[0], "arg": argsValue}
		if err := sjson.Encode(&buf, m); err != nil {
			return err
		}
	case Script:
		m := map[string]sjson.Value{"type": "script", "script": command}
		if err := sjson.Encode(&buf, m); err != nil {
			return err
		}
	default:
		return errors.New("invalid command type")
	}

	data := buf.Bytes()
	for n := 0; n < len(data); {
		num, err := con.ws.Write(buf.Bytes())
		if err != nil {
			return err
		}
		n += num
	}

	return nil
}

func (con *Console) SetDeadline(t time.Time) {
	con.ws.SetDeadline(t)
}

func (con *Console) Host() string {
	return con.host
}

func (con *Console) Close() {
	con.ws.Close()
}

func NewConsole(host string) (*Console, error) {
	h, p, err := net.SplitHostPort(host)
	if err != nil {
		h = host
		p = strconv.Itoa(DefaultPort)
	}

	addr := net.JoinHostPort(h, p)
	url := fmt.Sprintf("ws://%s", addr)
	ws, err := websocket.Dial(url, "", "http://"+h)
	if err != nil {
		return nil, err
	}

	con := &Console{sjson.NewLexer(ws), ws, addr}
	return con, nil
}
