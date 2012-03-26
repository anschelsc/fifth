package main

import (
	"fmt"
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

type fNull struct{}

func (_ fNull) run() {}

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

type fCap struct {
	place *fPush
}

func (c fCap) run() {
	c.place.d = dstack[len(dstack)-1]
	dstack = dstack[:len(dstack)-1]
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
			dstack = dstack[:l-1]
		}),
		".": fBuiltin(func() {
			l := len(dstack)
			fmt.Printf("%d ", dstack[l-1])
			dstack = dstack[:l-1]
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

func toThread(r []pFunc) fThread {
	names = append(names, make(map[string]Func)) // New scope
	compiled := make(fThread, len(r))
	for i, v := range r {
		compiled[i] = v.eval()
	}
	names = names[:len(names)-1] // End scope
	return compiled
}

func (l pLambdaFunc) eval() Func { return fPush{dFunc{toThread([]pFunc(l))}} }

func (n pNamedFunc) eval() Func { return toThread([]pFunc(n.inside)) }

func (c pCap) eval() Func {
	place := new(fPush)
	names[len(names)-1][string(c)] = place
	return fCap{place}
}

func process(a AST) {
	level := make(map[string]Func, len(a))
	for _, v := range a {
		level[v.name] = new(fThread)
	}
	names = append(names, level)
	for _, v := range a {
		*level[v.name].(*fThread) = v.eval().(fThread)
	}
}

func run() {
	rstack = rstack[:0]
	rstack = append(rstack, lookup("main"))
	for len(rstack) != 0 {
		next := rstack[len(rstack)-1]
		rstack = rstack[:len(rstack)-1]
		next.run()
	}
}
