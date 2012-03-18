package main

import (
	"unicode"
)

type token struct {
	k    kind
	data []rune
}

type kind uint8

const (
	kId kind = iota + 1
	kClose
	kFOpen // (named) function
	kLOpen // lambda
	kTOpen // type signature
	kCap
	kBS // backslash
)

func (t *token) String() string {
	if t == nil {
		return "LEXER ERROR"
	}
	switch t.k {
	case kId:
		return string(t.data)
	case kClose:
		return ")"
	case kFOpen:
		return "("
	case kLOpen:
		return "\\("
	case kTOpen:
		return ":("
	case kCap:
		return "@"
	case kBS:
		return "\\"
	}
	panic("Bad token")
}

type lexFunc func(rune, chan<- *token) lexFunc

func lex(source <-chan rune) <-chan *token {
	output := make(chan *token)
	go func() {
		f := lexMain
		for r := range source {
			f = f(r, output)
		}
		close(output)
	}()
	return output
}

func lexMain(r rune, output chan<- *token) lexFunc {
	if unicode.IsSpace(r) {
		return lexMain
	}
	switch r {
	case '@':
		output <- &token{kCap, nil}
		return lexMain
	case '(':
		output <- &token{kFOpen, nil}
		return lexMain
	case ')':
		output <- &token{kClose, nil}
		return lexMain
	case ':':
		return lexCol
	case '\\':
		return lexBS
	case '#':
		return lexComment
	}
	return lexBuf([]rune{r})
}

func lexBuf(buf []rune) lexFunc {
	return func(r rune, output chan<- *token) lexFunc {
		if unicode.IsSpace(r) {
			output <- &token{kId, buf}
			return lexMain
		}
		switch r {
		case '(', ')', '\\', '@', ':':
			output <- &token{kId, buf}
			return lexMain(r, output) // Note that we relex r
		}
		return lexBuf(append(buf, r))
	}
}

func lexCol(r rune, output chan<- *token) lexFunc {
	if unicode.IsSpace(r) {
		return lexCol
	}
	if r == '(' {
		output <- &token{kTOpen, nil}
		return lexMain
	}
	output <- nil
	return lexMain(r, output)
}

func lexBS(r rune, output chan<- *token) lexFunc {
	if unicode.IsSpace(r) {
		return lexBS
	}
	if r == '(' {
		output <- &token{kLOpen, nil}
		return lexMain
	}
	output <- &token{kBS, nil}
	return lexMain(r, output)
}

func lexComment(r rune, _ chan<- *token) lexFunc {
	if r == '\n' {
		return lexMain
	}
	return lexComment
}
