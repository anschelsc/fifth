package main

import (
	"bytes"
	"fmt"
)

type syntaxError struct {
	expected string
	got      string
}

func (e *syntaxError) Error() string {
	return fmt.Sprintf("Syntax error: expected %s, got %s.", e.expected, e.got)
}

type pFunc interface {
	String() string
}

type pSimpleFunc string

func (s pSimpleFunc) String() string { return string(s) }

type pPushFunc string

func (p pPushFunc) String() string { return string("\\" + p) }

type pLambdaFunc []pFunc

func (l pLambdaFunc) String() string {
	buf := new(bytes.Buffer)
	buf.Write([]byte("\\( "))
	for _, f := range l {
		buf.Write([]byte(f.String()))
		buf.WriteByte(' ')
	}
	buf.WriteByte(')')
	return buf.String()
}

type pNamedFunc struct {
	name   string
	inside []pFunc
}

func (p *pNamedFunc) String() string {
	buf := new(bytes.Buffer)
	buf.Write([]byte(p.name))
	buf.Write([]byte("( "))
	for _, f := range p.inside {
		buf.Write([]byte(f.String()))
		buf.WriteByte(' ')
	}
	buf.WriteByte(')')
	return buf.String()
}

type AST []*pNamedFunc

func parse(input <-chan *token) (AST, error) {
	ret := make(AST, 0)
	for t := range input {
		if t.k != kId {
			return ret, &syntaxError{"identifier", t.String()}
		}
		name := string(t.data)
		t, ok := <-input
		if !ok {
			return ret, &syntaxError{"(", "EOF"}
		}
		if t.k != kFOpen {
			return ret, &syntaxError{"(", t.String()}
		}
		inside, err := parseFunc(input)
		if err != nil {
			return ret, err
		}
		ret = append(ret, &pNamedFunc{name, inside})
	}
	return ret, nil
}

func parseFunc(input <-chan *token) ([]pFunc, error) {
	ret := make([]pFunc, 0)
	for {
		t, ok := <-input
		if !ok {
			return ret, &syntaxError{"identifier or \\( or )", "EOF"}
		}
		switch t.k {
		case kClose:
			return ret, nil
		case kId:
			ret = append(ret, pSimpleFunc(t.data))
		case kLOpen:
			inside, err := parseFunc(input)
			ret = append(ret, pLambdaFunc(inside)) // This way if there's some data it can be saved...
			if err != nil {
				return ret, err
			}
		case kBS:
			next, err := parseId(input)
			ret = append(ret, pPushFunc(next))
			if err != nil {
				return ret, err
			}
		default:
			return ret, &syntaxError{"identifier or \\( or )", t.String()}
		}
	}
	panic("unreachable")
}

func parseId(input <-chan *token) (string, error) {
	t, ok := <-input
	if !ok {
		return "", &syntaxError{"identifier", "EOF"}
	}
	if t.k != kId {
		return "", &syntaxError{"identifier", t.String()}
	}
	return string(t.data), nil
}
