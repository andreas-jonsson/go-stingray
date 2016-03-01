//line sjson.y:1

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

import __yyfmt__ "fmt"

//line sjson.y:19
//line sjson.y:22
type yySymType struct {
	yys int
	v   interface{}
}

const NUMBER = 57346
const STRING = 57347
const IDENTIFIER = 57348
const BOOLEAN = 57349
const NULL = 57350
const OBJECT_BEGIN = 57351
const OBJECT_END = 57352
const ARRAY_BEGIN = 57353
const ARRAY_END = 57354
const COMMA = 57355
const EQUAL = 57356

var yyToknames = [...]string{
	"$end",
	"error",
	"$unk",
	"NUMBER",
	"STRING",
	"IDENTIFIER",
	"BOOLEAN",
	"NULL",
	"OBJECT_BEGIN",
	"OBJECT_END",
	"ARRAY_BEGIN",
	"ARRAY_END",
	"COMMA",
	"EQUAL",
}
var yyStatenames = [...]string{}

const yyEofCode = 1
const yyErrCode = 2
const yyInitialStackSize = 16

//line sjson.y:63

//line yacctab:1
var yyExca = [...]int{
	-1, 1,
	1, -1,
	-2, 0,
}

const yyNprod = 18
const yyPrivate = 57344

var yyTokenNames []string
var yyStates []string

const yyLast = 36

var yyAct = [...]int{

	2, 4, 3, 13, 5, 6, 9, 22, 10, 16,
	21, 18, 4, 3, 17, 5, 6, 9, 8, 10,
	14, 15, 26, 27, 25, 28, 23, 24, 19, 14,
	15, 20, 12, 7, 11, 1,
}
var yyPact = [...]int{

	8, -1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000, 24,
	-3, -1000, 18, -1000, -4, -7, -1000, 14, -1000, -1000,
	15, 8, 8, -1000, 8, -1000, -1000, -1000, -1000,
}
var yyPgo = [...]int{

	0, 35, 0, 33, 32, 3, 18, 14,
}
var yyR1 = [...]int{

	0, 1, 3, 3, 4, 4, 5, 5, 6, 6,
	7, 7, 2, 2, 2, 2, 2, 2,
}
var yyR2 = [...]int{

	0, 1, 2, 3, 1, 3, 3, 3, 2, 3,
	1, 3, 1, 1, 1, 1, 1, 1,
}
var yyChk = [...]int{

	-1000, -1, -2, 5, 4, 7, 8, -3, -6, 9,
	11, 10, -4, -5, 5, 6, 12, -7, -2, 10,
	13, 14, 14, 12, 13, -5, -2, -2, -2,
}
var yyDef = [...]int{

	0, -2, 1, 12, 13, 14, 15, 16, 17, 0,
	0, 2, 0, 4, 0, 0, 8, 0, 10, 3,
	0, 0, 0, 9, 0, 5, 6, 7, 11,
}
var yyTok1 = [...]int{

	1,
}
var yyTok2 = [...]int{

	2, 3, 4, 5, 6, 7, 8, 9, 10, 11,
	12, 13, 14,
}
var yyTok3 = [...]int{
	0,
}

var yyErrorMessages = [...]struct {
	state int
	token int
	msg   string
}{}

//line yaccpar:1

/*	parser for yacc output	*/

var (
	yyDebug        = 0
	yyErrorVerbose = false
)

type yyLexer interface {
	Lex(lval *yySymType) int
	Error(s string)
}

type yyParser interface {
	Parse(yyLexer) int
	Lookahead() int
}

type yyParserImpl struct {
	lval  yySymType
	stack [yyInitialStackSize]yySymType
	char  int
}

func (p *yyParserImpl) Lookahead() int {
	return p.char
}

func yyNewParser() yyParser {
	return &yyParserImpl{}
}

const yyFlag = -1000

func yyTokname(c int) string {
	if c >= 1 && c-1 < len(yyToknames) {
		if yyToknames[c-1] != "" {
			return yyToknames[c-1]
		}
	}
	return __yyfmt__.Sprintf("tok-%v", c)
}

func yyStatname(s int) string {
	if s >= 0 && s < len(yyStatenames) {
		if yyStatenames[s] != "" {
			return yyStatenames[s]
		}
	}
	return __yyfmt__.Sprintf("state-%v", s)
}

func yyErrorMessage(state, lookAhead int) string {
	const TOKSTART = 4

	if !yyErrorVerbose {
		return "syntax error"
	}

	for _, e := range yyErrorMessages {
		if e.state == state && e.token == lookAhead {
			return "syntax error: " + e.msg
		}
	}

	res := "syntax error: unexpected " + yyTokname(lookAhead)

	// To match Bison, suggest at most four expected tokens.
	expected := make([]int, 0, 4)

	// Look for shiftable tokens.
	base := yyPact[state]
	for tok := TOKSTART; tok-1 < len(yyToknames); tok++ {
		if n := base + tok; n >= 0 && n < yyLast && yyChk[yyAct[n]] == tok {
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}
	}

	if yyDef[state] == -2 {
		i := 0
		for yyExca[i] != -1 || yyExca[i+1] != state {
			i += 2
		}

		// Look for tokens that we accept or reduce.
		for i += 2; yyExca[i] >= 0; i += 2 {
			tok := yyExca[i]
			if tok < TOKSTART || yyExca[i+1] == 0 {
				continue
			}
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}

		// If the default action is to accept or reduce, give up.
		if yyExca[i+1] != 0 {
			return res
		}
	}

	for i, tok := range expected {
		if i == 0 {
			res += ", expecting "
		} else {
			res += " or "
		}
		res += yyTokname(tok)
	}
	return res
}

func yylex1(lex yyLexer, lval *yySymType) (char, token int) {
	token = 0
	char = lex.Lex(lval)
	if char <= 0 {
		token = yyTok1[0]
		goto out
	}
	if char < len(yyTok1) {
		token = yyTok1[char]
		goto out
	}
	if char >= yyPrivate {
		if char < yyPrivate+len(yyTok2) {
			token = yyTok2[char-yyPrivate]
			goto out
		}
	}
	for i := 0; i < len(yyTok3); i += 2 {
		token = yyTok3[i+0]
		if token == char {
			token = yyTok3[i+1]
			goto out
		}
	}

out:
	if token == 0 {
		token = yyTok2[1] /* unknown char */
	}
	if yyDebug >= 3 {
		__yyfmt__.Printf("lex %s(%d)\n", yyTokname(token), uint(char))
	}
	return char, token
}

func yyParse(yylex yyLexer) int {
	return yyNewParser().Parse(yylex)
}

func (yyrcvr *yyParserImpl) Parse(yylex yyLexer) int {
	var yyn int
	var yyVAL yySymType
	var yyDollar []yySymType
	_ = yyDollar // silence set and not used
	yyS := yyrcvr.stack[:]

	Nerrs := 0   /* number of errors */
	Errflag := 0 /* error recovery flag */
	yystate := 0
	yyrcvr.char = -1
	yytoken := -1 // yyrcvr.char translated into internal numbering
	defer func() {
		// Make sure we report no lookahead when not parsing.
		yystate = -1
		yyrcvr.char = -1
		yytoken = -1
	}()
	yyp := -1
	goto yystack

ret0:
	return 0

ret1:
	return 1

yystack:
	/* put a state and value onto the stack */
	if yyDebug >= 4 {
		__yyfmt__.Printf("char %v in %v\n", yyTokname(yytoken), yyStatname(yystate))
	}

	yyp++
	if yyp >= len(yyS) {
		nyys := make([]yySymType, len(yyS)*2)
		copy(nyys, yyS)
		yyS = nyys
	}
	yyS[yyp] = yyVAL
	yyS[yyp].yys = yystate

yynewstate:
	yyn = yyPact[yystate]
	if yyn <= yyFlag {
		goto yydefault /* simple state */
	}
	if yyrcvr.char < 0 {
		yyrcvr.char, yytoken = yylex1(yylex, &yyrcvr.lval)
	}
	yyn += yytoken
	if yyn < 0 || yyn >= yyLast {
		goto yydefault
	}
	yyn = yyAct[yyn]
	if yyChk[yyn] == yytoken { /* valid shift */
		yyrcvr.char = -1
		yytoken = -1
		yyVAL = yyrcvr.lval
		yystate = yyn
		if Errflag > 0 {
			Errflag--
		}
		goto yystack
	}

yydefault:
	/* default state action */
	yyn = yyDef[yystate]
	if yyn == -2 {
		if yyrcvr.char < 0 {
			yyrcvr.char, yytoken = yylex1(yylex, &yyrcvr.lval)
		}

		/* look through exception table */
		xi := 0
		for {
			if yyExca[xi+0] == -1 && yyExca[xi+1] == yystate {
				break
			}
			xi += 2
		}
		for xi += 2; ; xi += 2 {
			yyn = yyExca[xi+0]
			if yyn < 0 || yyn == yytoken {
				break
			}
		}
		yyn = yyExca[xi+1]
		if yyn < 0 {
			goto ret0
		}
	}
	if yyn == 0 {
		/* error ... attempt to resume parsing */
		switch Errflag {
		case 0: /* brand new error */
			yylex.Error(yyErrorMessage(yystate, yytoken))
			Nerrs++
			if yyDebug >= 1 {
				__yyfmt__.Printf("%s", yyStatname(yystate))
				__yyfmt__.Printf(" saw %s\n", yyTokname(yytoken))
			}
			fallthrough

		case 1, 2: /* incompletely recovered error ... try again */
			Errflag = 3

			/* find a state where "error" is a legal shift action */
			for yyp >= 0 {
				yyn = yyPact[yyS[yyp].yys] + yyErrCode
				if yyn >= 0 && yyn < yyLast {
					yystate = yyAct[yyn] /* simulate a shift of "error" */
					if yyChk[yystate] == yyErrCode {
						goto yystack
					}
				}

				/* the current p has no shift on "error", pop stack */
				if yyDebug >= 2 {
					__yyfmt__.Printf("error recovery pops state %d\n", yyS[yyp].yys)
				}
				yyp--
			}
			/* there is no state on the stack with an error shift ... abort */
			goto ret1

		case 3: /* no shift yet; clobber input char */
			if yyDebug >= 2 {
				__yyfmt__.Printf("error recovery discards %s\n", yyTokname(yytoken))
			}
			if yytoken == yyEofCode {
				goto ret1
			}
			yyrcvr.char = -1
			yytoken = -1
			goto yynewstate /* try again in the same state */
		}
	}

	/* reduction by production yyn */
	if yyDebug >= 2 {
		__yyfmt__.Printf("reduce %v in:\n\t%v\n", yyn, yyStatname(yystate))
	}

	yynt := yyn
	yypt := yyp
	_ = yypt // guard against "declared and not used"

	yyp -= yyR2[yyn]
	// yyp is now the index of $0. Perform the default action. Iff the
	// reduced production is ε, $1 is possibly out of range.
	if yyp+1 >= len(yyS) {
		nyys := make([]yySymType, len(yyS)*2)
		copy(nyys, yyS)
		yyS = nyys
	}
	yyVAL = yyS[yyp+1]

	/* consult goto table to find next state */
	yyn = yyR1[yyn]
	yyg := yyPgo[yyn]
	yyj := yyg + yyS[yyp].yys + 1

	if yyj >= yyLast {
		yystate = yyAct[yyg]
	} else {
		yystate = yyAct[yyj]
		if yyChk[yystate] != -yyn {
			yystate = yyAct[yyg]
		}
	}
	// dummy call; replaced with literal code
	switch yynt {

	case 1:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sjson.y:33
		{
			yylex.(*Lexer).parseResult = yyDollar[1].v
			return 0
		}
	case 2:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line sjson.y:36
		{
			yyVAL.v = make(map[string]interface{})
		}
	case 3:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line sjson.y:37
		{
			yyVAL = yyDollar[2]
		}
	case 4:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sjson.y:40
		{
			m := make(map[string]interface{})
			p := yyDollar[1].v.([2]interface{})
			m[p[0].(string)] = p[1]
			yyVAL.v = m
		}
	case 5:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line sjson.y:41
		{
			m := yyDollar[1].v.(map[string]interface{})
			p := yyDollar[3].v.([2]interface{})
			m[p[0].(string)] = p[1]
		}
	case 6:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line sjson.y:44
		{
			yyVAL.v = [2]interface{}{yyDollar[1].v, yyDollar[3].v}
		}
	case 7:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line sjson.y:45
		{
			yyVAL.v = [2]interface{}{yyDollar[1].v, yyDollar[3].v}
		}
	case 8:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line sjson.y:48
		{
			yyVAL.v = make([]interface{}, 0)
		}
	case 9:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line sjson.y:49
		{
			yyVAL = yyDollar[2]
		}
	case 10:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sjson.y:52
		{
			s := make([]interface{}, 0, 1)
			s = append(s, yyDollar[1].v)
			yyVAL.v = s
		}
	case 11:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line sjson.y:53
		{
			yyVAL.v = append(yyDollar[1].v.([]interface{}), yyDollar[3].v)
		}
	case 17:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sjson.y:61
		{
			yyVAL = yyDollar[1]
		}
	}
	goto yystack /* stack new state and value */
}
