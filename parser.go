package main

import (
	"io"
	"errors"
	"fmt"
)

type chunk interface {
	WriteTo(io.Writer) error
}

type bangc struct{}

type capturec string

type identc string

type numc int

type closurec []chunk

func parse(tch <-chan token, ech <-chan error) ([]chunk, error) {
	ch, closed, err := parseInner(tch, ech)
	if err != nil {
		return nil, err
	}
	if closed {
		return nil, errors.New("Unmatched ')'.")
	}
	return ch, nil
}

func parseInner(tch <-chan token, ech <-chan error) ([]chunk, bool, error) {
	var ret []chunk
	for {
		t, ok := <-tch
		if !ok {
			err := <-ech
			return ret, false, err
		}
		switch t.tkind {
		case AT:
			t, ok = <-tch
			if !ok {
				err := <-ech
				if err == nil {
					err = errors.New("Expected IDENT after AT; got EOF.")
				}
				return ret, false, err
			}
			if t.tkind != IDENT {
				return ret, false, fmt.Errorf("Expected IDENT after AT; got %s.", t)
			}
			ret = append(ret, capturec(t.val))
		case BANG:
			ret = append(ret, bangc{})
		case OPEN:
			inner, closed, err := parseInner(tch, ech)
			ret = append(ret, closurec(inner))
			if err != nil {
				return ret, false, err
			}
			if !closed {
				return ret, false, errors.New("Unmatched '('")
			}
		case CLOSE:
			return ret, true, nil
		case NUM:
			ret = append(ret, numc(t.val[0]))
		case IDENT:
			ret = append(ret, identc(t.val))
		}
	}
	panic("unreachable")
}

func writeByte(w io.Writer, b byte) error {
	n, err := w.Write([]byte{b})
	if n != 1 {
		return err
	}
	return nil
}

func (_ bangc) WriteTo(w io.Writer) error {
	return writeByte(w, '!')
}

func (c capturec) WriteTo(w io.Writer) error {
	err := writeByte(w, '@')
	if err != nil {
		return err
	}
	_, err = io.WriteString(w, string(c))
	return err
}

func (i identc) WriteTo(w io.Writer) error {
	_, err := io.WriteString(w, string(i))
	return err
}

func (n numc) WriteTo(w io.Writer) error {
	_, err := fmt.Fprintf(w, "%d", int(n))
	return err
}

func (c closurec) WriteTo(w io.Writer) error {
	err := writeByte(w, '(')
	if err != nil {
		return err
	}
	first := true
	for _, ch := range c {
		if first {
			first = false
		} else {
			err := writeByte(w, ' ')
			if err != nil {
				return err
			}
		}
		err := ch.WriteTo(w)
		if err != nil {
			return err
		}
	}
	return writeByte(w, ')')
}
