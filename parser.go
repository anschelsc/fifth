package fifth

import (
	"errors"
	"fmt"
)

type chunk interface {
	unbound() []string
	eval(*world, []map[string]object) error
}

type bangc struct{}

type capturec string

type identc string

type numc int

type stringc []rune

type charc rune

type closurec struct {
	chunks   []chunk
	_unbound []string
}

type parseError struct {
	err    error
	closed bool
}

func parse(tch <-chan token, ech <-chan error) (<-chan chunk, <-chan error) {
	cChan, eChan := parseInner(tch, ech)
	outEChan := make(chan error)
	go func() {
		err := <-eChan
		if err.err != nil {
			outEChan <- err.err
			return
		}
		if err.closed {
			outEChan <- errors.New("Unmatched ')'.")
			return
		}
		outEChan <- nil
	}()
	return cChan, outEChan
}

func parseInner(tch <-chan token, ech <-chan error) (<-chan chunk, <-chan *parseError) {
	cChan := make(chan chunk)
	eChan := make(chan *parseError)
	go func() {
		for {
			t, ok := <-tch
			if !ok {
				close(cChan)
				err := <-ech
				eChan <- &parseError{err: err, closed: false}
				return
			}
			switch t.tkind {
			case AT:
				t, ok = <-tch
				if !ok {
					close(cChan)
					err := <-ech
					if err == nil {
						err = errors.New("Expected IDENT after AT; got EOF.")
					}
					eChan <- &parseError{err: err, closed: false}
					return
				}
				if t.tkind != IDENT {
					close(cChan)
					eChan <- &parseError{err: fmt.Errorf("Expected IDENT after AT; got %s.", t), closed: false}
					return
				}
				cChan <- capturec(t.val)
			case BANG:
				cChan <- bangc{}
			case OPEN:
				var inner []chunk
				inner_cChan, inner_eChan := parseInner(tch, ech)
				for ch := range inner_cChan {
					inner = append(inner, ch)
				}
				cChan <- &closurec{chunks: inner}
				err := <-inner_eChan
				if err.err != nil {
					close(cChan)
					eChan <- err
					return
				}
				if !err.closed {
					close(cChan)
					err.err = errors.New("Unmatched '('")
					eChan <- err
					return
				}
			case CLOSE:
				close(cChan)
				eChan <- &parseError{err: nil, closed: true}
				return
			case NUM:
				cChan <- numc(t.val[0])
			case IDENT:
				cChan <- identc(t.val)
			case STRING:
				cChan <- stringc(t.val)
			case CHAR:
				cChan <- charc(t.val[0])
			}
		}
	}()
	return cChan, eChan
}

func (_ bangc) unbound() []string { return []string{} }

func (_ capturec) unbound() []string { return []string{} }

func (i identc) unbound() []string { return []string{string(i)} }

func (_ numc) unbound() []string { return []string{} }

func (_ stringc) unbound() []string { return []string{} }

func (_ charc) unbound() []string { return []string{} }

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
