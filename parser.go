package main

import (
	"errors"
	"fmt"
)

type chunk interface {
	unbound() []string
	eval([]map[string]object) error
}

type bangc struct{}

type capturec string

type identc string

type numc int

type closurec struct {
	chunks []chunk
	_unbound []string
}

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
			ret = append(ret, &closurec{chunks: inner})
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

func (_ bangc) unbound() []string { return []string{} }

func (_ capturec) unbound() []string { return []string{} }

func (i identc) unbound() []string { return []string{string(i)} }

func (_ numc) unbound() []string { return []string{} }

func (c *closurec) unbound() []string {
	if c._unbound == nil {
		ret := make(map[string]struct{})
		bound := make(map[string]struct{})
		for _, ch := range c.chunks {
			if capt, ok := ch.(capturec); ok {
				bound[string(capt)] = struct{}{}
			} else {
				for _, s := range ch.unbound() {
					if _, ok = bound[s]; !ok {
						ret[s] = struct{}{}
					}
				}
			}
		}
		c._unbound = make([]string, 0, len(ret))
		for s := range ret {
			c._unbound = append(c._unbound, s)
		}
	}
	return c._unbound
}
