package main

type object interface {
}

var stack []object

func push(o object) { stack = append(stack, o) }

func pop() object {
	i := len(stack) - 1
	ret := stack[i]
	stack = stack[:i]
	return ret
}

type numo int

type stringo []rune

type charo rune

type function interface {
	object
	run([]map[string]object) error
}

type builtin func([]map[string]object) error

func (f builtin) run(context []map[string]object) error {
	return f(context)
}

type closure struct {
	todo     []chunk
	bindings map[string]object
}

func (c *closure) run(context []map[string]object) error {
	context = append(append(context, c.bindings), make(map[string]object))
	for _, ch := range c.todo {
		err := ch.eval(context)
		if err != nil {
			return err
		}
	}
	return nil
}
