package main

import (
	"io"
	"bufio"
	"unicode"
	"fmt"
)

type state func(chan<- token, rune) state

type token struct {
	tkind
	val []rune
}

type tkind byte

const (
	AT tkind = iota
	BANG
	OPEN
	CLOSE
	NUM
	IDENT
)

func (k tkind) String() string {
	return []string{
		AT: "AT",
		BANG: "BANG",
		OPEN: "OPEN",
		CLOSE: "CLOSE",
		NUM: "NUM",
		IDENT: "IDENT",
	}[k]
}

func (t token) String() string {
	switch t.tkind {
	case IDENT:
		return fmt.Sprintf("IDENT(%s)", string(t.val))
	case NUM:
		return fmt.Sprintf("NUM(%d)", t.val[0])
	}
	return t.tkind.String()
}

// token chan gets closed after all tokens are sent
// after which a single error is sent on the error chan
func lex(r io.Reader) (<-chan token, <-chan error) {
	tch := make(chan token)
	ech := make(chan error)
	go func() {
		br := bufio.NewReader(r)
		s := sstate
		ru, _, err := br.ReadRune()
		for ; err == nil; ru, _, err = br.ReadRune() {
			s = s(tch, ru)
		}
		close(tch)
		if err == io.EOF {
			err = nil
		}
		ech <- err
		close(ech)
	}()
	return tch, ech
}

func sstate(ch chan<- token, r rune) state {
	if unicode.IsSpace(r) {
		return sstate
	}
	switch r {
	case '@':
		ch <- token{tkind: AT}
		return sstate
	case '!':
		ch <- token{tkind: BANG}
		return sstate
	case '(':
		ch <- token{tkind: OPEN}
		return sstate
	case ')':
		ch <- token{tkind: CLOSE}
		return sstate
	}
	if '0' <= r && r <= '9' {
		return nstate(r - '0')
	}
	return istate([]rune{r})
}

func nstate(val rune) state {
	return func(ch chan<- token, r rune) state {
		if '0' <= r && r <= '9' {
			return nstate(val*10 + (r - '0'))
		}
		ch <- token{tkind: NUM, val: []rune{val}}
		return sstate(ch, r)
	}
}

func istate(val []rune) state {
	return func(ch chan<- token, r rune) state {
		if unicode.IsSpace(r) || r == '@' || r == '!' || r == '(' || r == ')' {
			ch <- token{tkind: IDENT, val: val}
			return sstate(ch, r)
		}
		return istate(append(val, r))
	}
}
