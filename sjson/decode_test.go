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
	"os"
	"testing"
)

func TestDecode(t *testing.T) {
	fp, err := os.Open("../testdata/test.json")
	if err != nil {
		t.Error(err)
	}
	defer fp.Close()

	lex := NewLexer(fp)

	_, err = Decode(lex)
	if err != nil {
		t.Error(err)
	}

	v, err := Decode(lex)
	if err != nil {
		t.Error(err)
	}

	if v.(string) != "next object" {
		t.Fail()
	}
}
