package fifth

import (
	"io"
)

type world struct {
	stack  []object
	input  io.Reader
	output io.Writer
}

func (w *world) run() error {
	cChan, eChan := parse(w.lex())
	return w.run_dyn(cChan, eChan, builtins)
}

func Run(input io.Reader, output io.Writer) error {
	w := &world{input: input, output: output}
	return w.run()
}
