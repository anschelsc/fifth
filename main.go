package main

import (
	"fmt"
	"os"
)

func main() {
	tch, ech := lex(os.Stdin)
	for t := range tch {
		fmt.Println(t)
	}
	if err := <-ech; err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
