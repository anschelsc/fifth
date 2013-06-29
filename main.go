package main

import (
	"fmt"
	"os"
)

func main() {
	cChan, eChan := parse(lex(os.Stdin))
	err := run_dyn(cChan, eChan, builtins)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
