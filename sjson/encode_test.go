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
	"bytes"
	"os"
	"reflect"
	"testing"
)

func TestEncode(t *testing.T) {
	fp, err := os.Open("../testdata/test.json")
	if err != nil {
		t.Error(err)
	}

	orginal, err := Decode(NewLexer(fp))
	if err != nil {
		t.Error(err)
	}
	fp.Close()

	var buf bytes.Buffer
	if err = Encode(&buf, orginal); err != nil {
		t.Error(err)
	}

	encoded, err := Decode(NewLexer(&buf))
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(orginal, encoded) {
		t.Error(err)
	}
}
