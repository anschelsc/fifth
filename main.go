package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	inCh := make(chan rune)
	go func() {
		r, _, err := in.ReadRune()
		for ; err == nil; r, _, err = in.ReadRune() {
			inCh <- r
		}
		if err != io.EOF {
			fmt.Fprintf(os.Stderr, "Error reading input: %s", err)
		}
		close(inCh)
	}()
	lexed := lex(inCh)
	ast, err := parse(lexed)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	process(ast)
	run()
}
