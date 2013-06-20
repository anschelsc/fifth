package main

import (
	"errors"
	"fmt"
)

var builtins = map[string]object{
	".": builtin(dot),
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
