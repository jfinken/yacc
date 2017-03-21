package main

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

const eof = -1

// Token is an alias for clarity but yacc requires
// the tokens to be defined in the grammar.
type Token int

// HocLex implements the HocLexer interface as defined
// by goyacc
type HocLex struct {
	src   string
	sidx  int
	start int
	width int
}

func (lxr *HocLex) lexOrig(lval *HocSymType) int {
	c := ' ' // rune
	for c == ' ' {
		if lxr.sidx == len(lxr.src) {
			return 0
		}
		c = rune(lxr.src[lxr.sidx])
		lxr.sidx++
	}

	if unicode.IsDigit(c) {
		f, _ := strconv.ParseFloat(string(c), 64)
		//fmt.Printf("int:   %d\n", int(c) - '0')
		fmt.Printf("float: %f\n", f)
		lval.val = f
		return DIGIT
	} else if unicode.IsLower(c) {
		lval.index = int(c) - 'a'
		return VAR
	}
	return int(c)
}

// Lex is the entry point for the lexer.  This func name signature and
// return type is expected by goyacc: implements HocLexer interface.
// yyparse calls Lex repeatedly for input tokens, returning a 0 signals
// "end of file" to yyparse.
func (lxr *HocLex) Lex(lval *HocSymType) int {

	tok := lxr.Next(lval)

	return int(tok)
}

// Error is a part of the HocLexer interface and prints syntax errors
func (lxr *HocLex) Error(s string) {
	fmt.Printf("syntax error: %s\n", s)
}

// Next parses each incoming src string to rune and returns
// a corresponding Token.  This is the core lexer entry point.
func (lxr *HocLex) Next(lval *HocSymType) Token {

	token := lxr.next()
	value := lxr.src[lxr.start:lxr.sidx]

	if token == DIGIT {
		f, _ := strconv.ParseFloat(value, 64)
		lval.val = f
	}
	//fmt.Printf("token: %v, value: %v, lval.val: %0.2f\n", token, value, lval.val)
	/*
		keyword := NIL
		if token == IDENTIFIER && lxr.peek() != ':' {
			keyword = Keyword(value)
		}
		return Item{token, keyword, value}
	*/
	return token
}

func (lxr *HocLex) next() Token {
	lxr.start = lxr.sidx
	c := lxr.read()
	// consume any white space
	lxr.whitespace(c)

	lxr.start = lxr.sidx
	c = lxr.read()
	switch c {
	//case eof:
	//	return Token(c)
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return lxr.number()
	case '.':
		if unicode.IsDigit(lxr.peek()) {
			return lxr.number()
		}
	default:
		if unicode.IsLetter(c) || c == '_' {
			return lxr.identifier()
		}
	}
	return Token(c)
}

func (lxr *HocLex) read() rune {
	if lxr.sidx >= len(lxr.src) {
		lxr.width = 0
		return eof
	}
	c, w := utf8.DecodeRuneInString(lxr.src[lxr.sidx:])
	lxr.sidx += w
	lxr.width = w
	return c
}
func (lxr *HocLex) peek() rune {
	c := lxr.read()
	lxr.backup()
	return c
}
func (lxr *HocLex) whitespace(c rune) {
	// just consume white space and tab
	for ; isSpace(c); c = lxr.read() {
		/* handled by the grammar
		if c == '\n' || c == '\r' {
			return NEWLINE
		}
		*/
	}
	lxr.backup()
}
func (lxr *HocLex) identifier() Token {
	lxr.matchWhile(isIdentChar)
	/*
		if !lxr.match('?') {
			lxr.match('!')
		}
	*/
	return IDENT
}
func (lxr *HocLex) number() Token {
	// No hex support
	digits := "0123456789"
	lxr.matchRunOf(digits)
	if lxr.match('.') {
		lxr.matchRunOf(digits)
	}
	return DIGIT
}
func (lxr *HocLex) match(c rune) bool {
	if c == lxr.read() {
		return true
	}
	lxr.backup()
	return false
}
func (lxr *HocLex) matchRunOf(valid string) {
	for strings.ContainsRune(valid, lxr.read()) {
	}
	lxr.backup()
}
func (lxr *HocLex) matchWhile(f func(c rune) bool) {
	for c := lxr.read(); f(c); c = lxr.read() {
	}
	lxr.backup()
}
func (lxr *HocLex) backup() {
	lxr.sidx -= lxr.width
}
func isSpace(c rune) bool {
	//return c == ' ' || c == '\t' || c == '\r' || c == '\n'
	// newline handled by grammar
	return c == ' ' || c == '\t'
}
func isIdentChar(r rune) bool {
	return r == '_' || unicode.IsLetter(r) || unicode.IsDigit(r)
}
