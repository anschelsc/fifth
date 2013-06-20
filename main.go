package main

import (
	"fmt"
	"os"
)

func main() {
	program, err := parse(lex(os.Stdin))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	cl := &closure{todo: program, bindings: builtins}
	err = cl.run(nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
