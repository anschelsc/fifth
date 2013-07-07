package fifth

type object interface {
}

func (w *world) push(o object) { w.stack = append(w.stack, o) }

func (w *world) pop() object {
	i := len(w.stack) - 1
	ret := w.stack[i]
	w.stack = w.stack[:i]
	return ret
}

type numo int

type stringo []rune

type charo rune

type function interface {
	object
	run(*world, []map[string]object) error
}

type builtin func(*world, []map[string]object) error

func (f builtin) run(w *world, context []map[string]object) error {
	return f(w, context)
}

type closure struct {
	todo     []chunk
	bindings map[string]object
}

func (c *closure) run(w *world, context []map[string]object) error {
	context = append(append(context, c.bindings), make(map[string]object))
	for _, ch := range c.todo {
		err := ch.eval(w, context)
		if err != nil {
			return err
		}
	}
	return nil
}

func (w *world) run_dyn(cChan <-chan chunk, eChan <-chan error, bindings map[string]object) error {
	context := []map[string]object{bindings, make(map[string]object)}
	for ch := range cChan {
		err := ch.eval(w, context)
		if err != nil {
			return err
		}
	}
	return <-eChan
}
