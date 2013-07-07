package main

import (
	"fmt"
	"os"

	"fifth"
)

func main() {
	err := fifth.Run(os.Stdin, os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
