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
	"encoding/binary"

	"github.com/andreas-jonsson/go-stingray/sjson"
)

type (
	Profiler struct {
		con *Console
	}

	ProfilerEvent struct {
		Type     uint32
		Name     uint64
		ThreadID uint32
		CoreID   uint32

		Parent      int32
		FirstChild  int32
		LastChild   int32
		PrevSibling int32
		NextSibling int32

		Time    float64
		Elapsed float64
		Count   uint32
	}
)

func (prof *Profiler) Pull(pe *ProfilerEvent) (sjson.Value, *ProfilerEvent, error) {
	for {
		val, data, err := prof.con.Receive()
		if err != nil {
			return nil, nil, err
		}

		if len(data) > 0 {
			if pe != nil {
				if err := binary.Read(bytes.NewReader(data), binary.BigEndian, pe); err != nil {
					return nil, nil, err
				}
			}
			return pe, nil, nil
		}

		if m, ok := val.(map[string]sjson.Value); ok {
			if ty, ok := m["type"]; ok {
				switch ty {
				case "profiler_strings", "profiler_threads":
					return val, nil, nil
				}
			}
		}
	}
}

func NewProfiler(con *Console) *Profiler {
	return &Profiler{con}
}
