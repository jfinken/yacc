package main

import (
	"fmt"
	"strconv"
	"unicode"
	"unicode/utf8"
)

const eof = -1

type Token int

const (
	NIL Token = iota
	EOF
	ERROR
)

// HocLex implements the HocLexer interface as defined
// by goyacc
type HocLex struct {
	src   string
	sidx  int
	start int
	width int
}

// Next parses each incoming src string to rune and returns
// a corresponding Token.  This is the core lexer entry point.
func (lxr *HocLex) Next() Token {

	token := lxr.next()

	value := lxr.src[lxr.start:lxr.sidx]
	fmt.Printf("%s\n", value)
	/*
		keyword := NIL
		if token == IDENTIFIER && lxr.peek() != ':' {
			keyword = Keyword(value)
		}
		return Item{token, keyword, value}
	*/
	return token
}

// Lex is the entry point for the lexer.  This func name signature and
// return type is expected by goyacc: implements HocLexer interface
func (lxr *HocLex) Lex(lval *HocSymType) int {
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

// Error is a part of the HocLexer interface and prints syntax errors
func (lxr *HocLex) Error(s string) {
	fmt.Printf("syntax error: %s\n", s)
}

func (lxr *HocLex) next() Token {
	lxr.start = lxr.sidx
	c := lxr.read()
	switch c {
	case eof:
		return EOF
	}
	// ...
	return ERROR
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
