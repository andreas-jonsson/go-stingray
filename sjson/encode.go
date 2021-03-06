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

package sjson

import (
	"errors"
	"fmt"
	"io"
	"strconv"
)

//Value represents a SJSON value.
type Value interface{}

//Encode encodes a SJSON value to the writer.
func Encode(writer io.Writer, v Value) error {
	return encodeValue(writer, v)
}

func encodeValue(writer io.Writer, v Value) error {
	var err error
	switch v.(type) {
	case nil:
		_, err = fmt.Fprint(writer, "null")
	case int, int8, int16, int32, int64,
		uint8, uint16, uint32, uint64,
		float32, float64, bool:
		_, err = fmt.Fprintf(writer, "%v", v)
	case string:
		_, err = fmt.Fprint(writer, strconv.Quote(v.(string)))
	case []Value:
		if _, err = fmt.Fprint(writer, "["); err != nil {
			return err
		}
		for i, val := range v.([]Value) {
			if i > 0 {
				if _, err = fmt.Fprint(writer, ","); err != nil {
					return err
				}
			}
			if err := encodeValue(writer, val); err != nil {
				return err
			}
		}
		_, err = fmt.Fprint(writer, "]")
	case map[string]Value:
		if _, err = fmt.Fprint(writer, "{"); err != nil {
			return err
		}

		i := 0
		for k, val := range v.(map[string]Value) {
			if i > 0 {
				if _, err = fmt.Fprint(writer, ","); err != nil {
					return err
				}
			}
			if _, err = fmt.Fprintf(writer, "\"%v\"=", k); err != nil {
				return err
			}
			if err := encodeValue(writer, val); err != nil {
				return err
			}
			i++
		}
		_, err = fmt.Fprint(writer, "}")
	default:
		return errors.New("invalid type")
	}
	return err
}
