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

func marshalMessage(v interface{}) ([]byte, byte, error) {
	return v.([]byte), websocket.TextFrame, nil
}

func unmarshalMessage(data []byte, ty byte, v interface{}) error {
	switch ty {
	case websocket.BinaryFrame:
		return nil
	case websocket.TextFrame:
		var err error
		lex := sjson.NewLexer(bytes.NewReader(data))

		val, err := sjson.Decode(lex)
		if err != nil {
			return err
		}

		defer func() {
			if r := recover(); r != nil {
				err = errors.New("invalid message")
			}
		}()

		msg := v.(*Message)
		m := val.(map[string]sjson.Value)

		if m["type"].(string) == "message" {
			msg.System = m["system"].(string)
			msg.Message = m["message"].(string)
			msg.MessageType = m["message_type"].(string)
			msg.Level = m["level"].(string)
			return err
		}

		return err
	default:
		return errors.New("unknown message")
	}
}

func unmarshalFrameData(data []byte, ty byte, v interface{}) error {
	fd := v.(*frameData)

	switch ty {
	case websocket.BinaryFrame:
		fd.data = data
	case websocket.TextFrame:
		lex := sjson.NewLexer(bytes.NewReader(data))
		val, err := sjson.Decode(lex)
		if err != nil {
			return err
		}
		fd.obj = val
	default:
		return errors.New("unknown message")
	}

	return nil
}

var (
	consoleMessageCodec   = websocket.Codec{Marshal: marshalMessage, Unmarshal: unmarshalMessage}
	consoleFrameDataCodec = websocket.Codec{Marshal: nil, Unmarshal: unmarshalFrameData}
)

type frameData struct {
	data []byte
	obj  sjson.Value
}

type Message struct {
	System,
	Level,
	MessageType,
	Message string
}

func (m Message) String() string {
	if m.System == "" {
		return "[?] " + m.Message
	} else {
		return fmt.Sprintf("[%s] %s", m.System, m.Message)
	}
}

type Console struct {
	lex  *sjson.Lexer
	ws   *websocket.Conn
	host string
}

func (con *Console) Receive() (sjson.Value, []byte, error) {
	fd := frameData{}
	err := consoleFrameDataCodec.Receive(con.ws, &fd)
	return fd.obj, fd.data, err
}

func (con *Console) ReceiveMessage() (Message, error) {
	var (
		err error
		msg Message
	)

	defer func() {
		if r := recover(); r != nil {
			err = errors.New("invalid message")
		}
	}()

	for msg.MessageType == "" {
		if err := consoleMessageCodec.Receive(con.ws, &msg); err != nil {
			return msg, err
		}
	}

	return msg, nil
}

func (con *Console) SendCommand(ty CommandType, command string) error {
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

	return consoleMessageCodec.Send(con.ws, buf.Bytes())
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

func NewConsole(host, protocol string) (*Console, error) {
	h, p, err := net.SplitHostPort(host)
	if err != nil {
		h = host
		p = strconv.Itoa(DefaultPort)
	}

	addr := net.JoinHostPort(h, p)
	url := fmt.Sprintf("ws://%s/%s", addr, protocol)

	ws, err := websocket.Dial(url, "", "http://"+h)
	if err != nil {
		return nil, err
	}

	con := &Console{sjson.NewLexer(ws), ws, addr}
	return con, nil
}
