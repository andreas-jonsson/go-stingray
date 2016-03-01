%{/*
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

%}

%union{
	v interface{}
}

%token _NUMBER _STRING _IDENTIFIER _BOOLEAN _NULL

%left _OBJECT_BEGIN _OBJECT_END _ARRAY_BEGIN _ARRAY_END
%left _COMMA _EQUAL

%%

START 		: _VALUE 								{ yylex.(*Lexer).parseResult = $1.v; return 0 }
			;

_OBJECT 	: _OBJECT_BEGIN _OBJECT_END 			{ $$.v = make(map[string]interface{}) }
			| _OBJECT_BEGIN _MEMBERS _OBJECT_END 	{ $$ = $2 }
			;

_MEMBERS 	: _PAIR 								{ m := make(map[string]interface{}); p := $1.v.([2]interface{}); m[p[0].(string)] = p[1]; $$.v = m }
			| _MEMBERS _COMMA _PAIR 				{ m := $1.v.(map[string]interface{}); p := $3.v.([2]interface{}); m[p[0].(string)] = p[1] }
			;

_PAIR 		: _STRING _EQUAL _VALUE 				{ $$.v = [2]interface{}{$1.v, $3.v} }
			| _IDENTIFIER _EQUAL _VALUE 			{ $$.v = [2]interface{}{$1.v, $3.v} }
			;

_ARRAY 		: _ARRAY_BEGIN _ARRAY_END  				{ $$.v = make([]interface{}, 0) }
			| _ARRAY_BEGIN _ELEMENTS _ARRAY_END 	{ $$ = $2 }
			;

_ELEMENTS	: _VALUE 								{ s := make([]interface{}, 0, 1); s = append(s, $1.v); $$.v = s }
			| _ELEMENTS _COMMA _VALUE 				{ $$.v = append($1.v.([]interface{}), $3.v) }
			;

_VALUE		: _STRING
			| _NUMBER
			| _BOOLEAN
			| _NULL
			| _OBJECT
			| _ARRAY 								{ $$ = $1 }
			;
%%