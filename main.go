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
	for _, ch := range program {
		ch.WriteTo(os.Stdout)
		fmt.Printf(" ")
	}
}
