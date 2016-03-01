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

%token NUMBER STRING IDENTIFIER BOOLEAN NULL

%left OBJECT_BEGIN OBJECT_END ARRAY_BEGIN ARRAY_END
%left COMMA EQUAL

%%

START 		: VALUE 							{ yylex.(*Lexer).parseResult = $1.v; return 0 }
			;

OBJECT 		: OBJECT_BEGIN OBJECT_END 			{ $$.v = make(map[string]interface{}) }
			| OBJECT_BEGIN MEMBERS OBJECT_END 	{ $$ = $2 }
			;

MEMBERS 	: PAIR 								{ m := make(map[string]interface{}); p := $1.v.([2]interface{}); m[p[0].(string)] = p[1]; $$.v = m }
			| MEMBERS COMMA PAIR 				{ m := $1.v.(map[string]interface{}); p := $3.v.([2]interface{}); m[p[0].(string)] = p[1] }
			;

PAIR 		: STRING EQUAL VALUE 				{ $$.v = [2]interface{}{$1.v, $3.v} }
			| IDENTIFIER EQUAL VALUE 			{ $$.v = [2]interface{}{$1.v, $3.v} }
			;

ARRAY 		: ARRAY_BEGIN ARRAY_END  			{ $$.v = make([]interface{}, 0) }
			| ARRAY_BEGIN ELEMENTS ARRAY_END 	{ $$ = $2 }
			;

ELEMENTS	: VALUE 							{ s := make([]interface{}, 0, 1); s = append(s, $1.v); $$.v = s }
			| ELEMENTS COMMA VALUE 				{ $$.v = append($1.v.([]interface{}), $3.v) }
			;

VALUE		: STRING
			| NUMBER
			| BOOLEAN
			| NULL
			| OBJECT
			| ARRAY 							{ $$ = $1 }
			;
%%
