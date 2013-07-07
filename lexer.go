package fifth

import (
	"bufio"
	"fmt"
	"io"
	"unicode"
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
	STRING
	CHAR
)

func (k tkind) String() string {
	return []string{
		AT:     "AT",
		BANG:   "BANG",
		OPEN:   "OPEN",
		CLOSE:  "CLOSE",
		NUM:    "NUM",
		IDENT:  "IDENT",
		STRING: "STRING",
		CHAR:   "CHAR",
	}[k]
}

func (t token) String() string {
	switch t.tkind {
	case IDENT:
		return fmt.Sprintf("IDENT(%s)", string(t.val))
	case NUM:
		return fmt.Sprintf("NUM(%d)", t.val[0])
	case STRING:
		return fmt.Sprintf("STRING(%q)", string(t.val))
	case CHAR:
		return fmt.Sprintf("CHAR(%c)", t.val[0])
	}
	return t.tkind.String()
}

// token chan gets closed after all tokens are sent
// after which a single error is sent on the error chan
func (w *world) lex() (<-chan token, <-chan error) {
	tch := make(chan token)
	ech := make(chan error)
	go func() {
		br := bufio.NewReader(w.input)
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
	case '{':
		return cstate(sstate)
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
	case '"':
		return qstate([]rune{})
	case '\'':
		return chstate
	}
	if '0' <= r && r <= '9' {
		return nstate(r - '0')
	}
	return istate([]rune{r})
}

func cstate(s state) state {
	return func(_ chan<- token, r rune) state {
		if r == '}' {
			return s
		}
		return cstate(s)
	}
}

func nstate(val rune) state {
	return func(ch chan<- token, r rune) state {
		if r == '{' {
			return cstate(nstate(val))
		}
		if '0' <= r && r <= '9' {
			return nstate(val*10 + (r - '0'))
		}
		ch <- token{tkind: NUM, val: []rune{val}}
		return sstate(ch, r)
	}
}

func istate(val []rune) state {
	return func(ch chan<- token, r rune) state {
		if r == '{' {
			return cstate(istate(val))
		}
		if unicode.IsSpace(r) || r == '@' || r == '!' || r == '(' || r == ')' {
			ch <- token{tkind: IDENT, val: val}
			return sstate(ch, r)
		}
		return istate(append(val, r))
	}
}

func qstate(val []rune) state {
	return func(ch chan<- token, r rune) state {
		switch r {
		case '"':
			ch <- token{tkind: STRING, val: val}
			return sstate
		case '\\':
			return bstate(val)
		}
		return qstate(append(val, r))
	}
}

func bstate(val []rune) state {
	return func(ch chan<- token, r rune) state {
		switch r {
		case '\\', '"':
			return qstate(append(val, r))
		case 'n':
			return qstate(append(val, '\n'))
		case 't':
			return qstate(append(val, '\t'))
		case '\n':
			return qstate(val)
		}
		return qstate(append(val, unicode.ReplacementChar))
	}
}

func chstate(ch chan<- token, r rune) state {
	ch <- token{tkind: CHAR, val: []rune{r}}
	return sstate
}
