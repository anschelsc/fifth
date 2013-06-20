package main

import (
	"fmt"
	"os"
)

func main() {
	program, err := parse(lex(os.Stdin))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	cl := &closure{todo: program, bindings: builtins}
	cl.run(nil)
}
