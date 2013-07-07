package fifth

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

func (_ bangc) eval(w *world, context []map[string]object) error {
	if len(w.stack) == 0 {
		return emptyStack
	}
	f, ok := w.pop().(function)
	if !ok {
		return errors.New("Tried to run a non-function.")
	}
	return f.run(w, context)
}

func (c capturec) eval(w *world, context []map[string]object) error {
	if len(w.stack) == 0 {
		return emptyStack
	}
	context[len(context)-1][string(c)] = w.pop()
	return nil
}

func (id identc) eval(w *world, context []map[string]object) error {
	ob, err := search(context, string(id))
	if err != nil {
		return err
	}
	w.push(ob)
	return nil
}

func (n numc) eval(w *world, _ []map[string]object) error {
	w.push(numo(n))
	return nil
}

func (s stringc) eval(w *world, _ []map[string]object) error {
	w.push(stringo(s))
	return nil
}

func (c charc) eval(w *world, _ []map[string]object) error {
	w.push(charo(c))
	return nil
}

func (c *closurec) eval(w *world, context []map[string]object) error {
	bindings := make(map[string]object)
	for _, id := range c.unbound() {
		ob, err := search(context, id)
		if err != nil {
			return err
		}
		bindings[id] = ob
	}
	w.push(&closure{todo: c.chunks, bindings: bindings})
	return nil
}
