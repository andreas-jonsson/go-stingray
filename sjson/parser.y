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
	v Value
}

%token _NUMBER _STRING _IDENTIFIER _BOOLEAN _NULL

%left _OBJECT_BEGIN _OBJECT_END _ARRAY_BEGIN _ARRAY_END
%left _COMMA _EQUAL _COLON

%%

START 		: _VALUE 								{ yylex.(*Lexer).parseResult = $1.v; return 0 }
			;

_OBJECT 	: _OBJECT_BEGIN _OBJECT_END 			{ $$.v = make(map[string]Value) }
			| _OBJECT_BEGIN _MEMBERS _OBJECT_END 	{ $$ = $2 }
			;

_MEMBERS 	: _PAIR 								{ m := make(map[string]Value); p := $1.v.([2]Value); m[p[0].(string)] = p[1]; $$.v = m }
			| _MEMBERS _PAIR						{ m := $1.v.(map[string]Value); p := $2.v.([2]Value); m[p[0].(string)] = p[1] }
			| _MEMBERS _COMMA _PAIR 				{ m := $1.v.(map[string]Value); p := $3.v.([2]Value); m[p[0].(string)] = p[1] }
			;

_PAIR 		: _STRING _SEP _VALUE 					{ $$.v = [2]Value{$1.v, $3.v} }
			| _IDENTIFIER _SEP _VALUE 				{ $$.v = [2]Value{$1.v, $3.v} }
			;

_SEP 		: _EQUAL
			| _COLON 								{}
			;

_ARRAY 		: _ARRAY_BEGIN _ARRAY_END  				{ $$.v = make([]Value, 0) }
			| _ARRAY_BEGIN _ELEMENTS _ARRAY_END 	{ $$ = $2 }
			;

_ELEMENTS	: _VALUE 								{ s := make([]Value, 0, 1); s = append(s, $1.v); $$.v = s }
			| _ELEMENTS _VALUE						{ $$.v = append($1.v.([]Value), $2.v) }
			| _ELEMENTS _COMMA _VALUE 				{ $$.v = append($1.v.([]Value), $3.v) }
			;

_VALUE		: _STRING
			| _NUMBER
			| _BOOLEAN
			| _NULL
			| _OBJECT
			| _ARRAY 								{ $$ = $1 }
			;
%%
