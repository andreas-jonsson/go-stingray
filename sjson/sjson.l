/[ \t\n]/               { }
/[0-9]+/                { i, _ := strconv.Atoi(yylex.Text()); lval.v = float64(i); return NUMBER }
/[0-9]+\.[0-9]*/        { lval.v, _ = strconv.ParseFloat(yylex.Text(), 64); return NUMBER }
/true|false/            { lval.v = (yylex.Text() == "true"); return BOOLEAN }
/null/                  { lval.v = nil; return NULL }
/[a-z][a-z0-9]*/        { lval.v = yylex.Text(); return IDENTIFIER }
/\"([^\\\"]|\\.)*\"/    { t := yylex.Text(); lval.v = t[1:len(t) - 1]; return STRING }

/=/     { return EQUAL }
/,/     { return COMMA }
/\[/    { return ARRAY_BEGIN }
/\]/    { return ARRAY_END }
/\{/    { return OBJECT_BEGIN }
/\}/    { return OBJECT_END }

/./     { panic(fmt.Errorf("unrecognized character: '%v'", yylex.Text())) }
//

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
    "strconv"
    "fmt"
)

type Decoder struct {
    lex *Lexer
}

func (dec *Decoder) Decode() (interface{}, error) {
    var err error
    v := func() interface{} {
        defer func() {
            if r := recover(); r != nil {
                err = fmt.Errorf("%v [%v:%v]", r, dec.lex.Line(), dec.lex.Column())
            }
        }()
        yyParse(dec.lex)
        return dec.lex.parseResult
    }()
    return v, err
}

func NewDecoder(reader io.Reader) *Decoder {
    return &Decoder{NewLexer(reader)}
}
