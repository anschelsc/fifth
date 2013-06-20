package main

import (
	"errors"
	"fmt"
)

var (
	emptyStack error = errors.New("Tried to pop from an empty stack.")
)

func search(context []map[string]object, ident string) (object, error) {
	for i := len(context) - 1; i >= 0; i-- {
		if ob, ok := context[i][ident]; ok {
			return ob, nil
		}
	}
	return nil, fmt.Errorf("Unbound identifier: %s.", ident)
}

func (_ bangc) eval(context []map[string]object) error {
	if len(stack) == 0 {
		return emptyStack
	}
	f, ok := pop().(function)
	if !ok {
		return errors.New("Tried to run a non-function.")
	}
	return f.run(context)
}

func (c capturec) eval(context []map[string]object) error {
	if len(stack) == 0 {
		return emptyStack
	}
	context[len(context) - 1][string(c)] = pop()
	return nil
}

func (id identc) eval(context []map[string]object) error {
	ob, err := search(context, string(id))
	if err != nil {
		return err
	}
	push(ob)
	return nil
}

func (n numc) eval(_ []map[string]object) error {
	push(numo(n))
	return nil
}

func (c *closurec) eval(context []map[string]object) error {
	bindings := make(map[string]object)
	for _, id := range c.unbound() {
		ob, err := search(context, id)
		if err != nil {
			return err
		}
		bindings[id] = ob
	}
	push(&closure{todo: c.chunks, bindings: bindings})
	return nil
}
