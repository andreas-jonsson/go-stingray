/[ \t\n]/                                   { }
/\-?[0-9]+/                                 { i, _ := strconv.Atoi(yylex.Text()); lval.v = float64(i); return _NUMBER }
/[-+]?[0-9]*\.?[0-9]+([eE][-+]?[0-9]+)?/    { lval.v, _ = strconv.ParseFloat(yylex.Text(), 64); return _NUMBER }
/true|false/                                { lval.v = (yylex.Text() == "true"); return _BOOLEAN }
/null/                                      { lval.v = nil; return _NULL }
/[a-z][a-z0-9]*/                            { lval.v = yylex.Text(); return _IDENTIFIER }
/\"([^\\\"]|\\.)*\"/                        { t := yylex.Text(); lval.v = t[1:len(t) - 1]; return _STRING }

/=/     { return _EQUAL }
/,/     { return _COMMA }
/\[/    { return _ARRAY_BEGIN }
/\]/    { return _ARRAY_END }
/\{/    { return _OBJECT_BEGIN }
/\}/    { return _OBJECT_END }

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

//Decode decodes a SJSON value from the stream.
func Decode(lex *Lexer) (Value, error) {
    var err error
    v := func() interface{} {
        defer func() {
            if r := recover(); r != nil {
                err = fmt.Errorf("%v [%v:%v]", r, lex.Line(), lex.Column())
            }
        }()
        yyParse(lex)
        return lex.parseResult
    }()
    return v, err
}
