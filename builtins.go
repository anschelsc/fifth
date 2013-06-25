package main

import (
	"errors"
	"fmt"
)

var builtins = map[string]object{
	"fail":  builtin(fail),
	".":     builtin(dot),
	"+":     builtin(plus),
	"_":     builtin(negate),
	"zero?": builtin(isZero),
	"%":     builtin(mod),
}

func fail(_ []map[string]object) error {
	return errors.New("Call to fail.")
}

func dot(_ []map[string]object) error {
	if len(stack) == 0 {
		return emptyStack
	}
	n, ok := pop().(numo)
	if !ok {
		return errors.New("Only numbers can be dot-printed.")
	}
	_, err := fmt.Printf("%d\n", int(n))
	return err
}

func plus(_ []map[string]object) error {
	if len(stack) < 2 {
		return emptyStack
	}
	n1, ok := pop().(numo)
	if !ok {
		return errors.New("Only numbers can be added.")
	}
	n2, ok := pop().(numo)
	if !ok {
		return errors.New("Only numbers can be added.")
	}
	push(n1 + n2)
	return nil
}

func negate(_ []map[string]object) error {
	if len(stack) == 0 {
		return emptyStack
	}
	n, ok := pop().(numo)
	if !ok {
		return errors.New("Only numbers can be negated.")
	}
	push(-n)
	return nil
}

func isZero(_ []map[string]object) error {
	if len(stack) == 0 {
		return emptyStack
	}
	n, ok := pop().(numo)
	if !ok {
		return errors.New("Only numbers can be checked for zeroness.")
	}
	if n == 0 {
		push(innerTrue)
	} else {
		push(innerFalse)
	}
	return nil
}

func mod(_ []map[string]object) error {
	if len(stack) < 2 {
		return emptyStack
	}
	n1, ok := pop().(numo)
	if !ok {
		return errors.New("Mod is for numbers.")
	}
	n2, ok := pop().(numo)
	if !ok {
		return errors.New("Mod is for numbers.")
	}
	push(n2 % n1)
	return nil
}

var innerTrue = builtin(func(context []map[string]object) error {
	if len(stack) < 2 {
		return emptyStack
	}
	pop()
	return bangc.eval(bangc{}, context)
})

var innerFalse = builtin(func(context []map[string]object) error {
	if len(stack) < 2 {
		return emptyStack
	}
	f := pop()
	pop()
	push(f)
	return bangc.eval(bangc{}, context)
})
