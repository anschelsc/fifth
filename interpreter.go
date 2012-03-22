package main

import (
	"strconv"
)

type Func interface {
	run()
}

type Data interface {
	String() string
}

var (
	rstack []Func
	dstack []Data
)

type fBuiltin func()

func (f fBuiltin) run() { f() }

type fPush struct {
	d Data
}

func (p fPush) run() {
	dstack = append(dstack, p.d)
}

type fThread []Func

func (t fThread) run() {
	for i := len(t) - 1; i >= 0; i-- {
		rstack = append(rstack, t[i])
	}
}

type dNum int

func (n dNum) String() string { return strconv.Itoa(int(n)) }

type dFunc struct {
	f Func
}

func (f dFunc) String() string { return "Function" }

var (
	builtins = map[string]Func{
		"do": fBuiltin(func() {
			rstack = append(rstack, dstack[len(dstack)-1].(dFunc).f)
			dstack = dstack[:len(dstack)-1]
		}),
		"+": fBuiltin(func() {
			l := len(dstack)
			dstack[l-2] = dstack[l-1].(dNum) + dstack[l-2].(dNum)
			dstack = dstack[:len(dstack)-1]
		}),
	}
	names = []map[string]Func{builtins}
)

func lookup(key string) Func {
	for i := len(names) - 1; i >= 0; i-- {
		if f, ok := names[i][key]; ok {
			return f
		}
	}
	val, _ := strconv.Atoi(key)
	return fPush{dNum(val)}
}

func (s pSimpleFunc) eval() Func { return lookup(string(s)) }

func (p pPushFunc) eval() Func { return fPush{dFunc{lookup(string(p))}} }

func (l pLambdaFunc) eval() Func {
	compiled := make(fThread, len(l))
	for i, v := range l {
		compiled[i] = v.eval()
	}
	return compiled
}

func (n pNamedFunc) eval() Func { return pLambdaFunc(n.inside).eval() }
