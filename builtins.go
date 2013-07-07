package fifth

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
	"c.": builtin(cDot),
	"c=": builtin(cEq),
	"sIndex": builtin(sIndex),
	"sLen": builtin(sLen),
}

func fail(_ *world, _ []map[string]object) error {
	return errors.New("Call to fail.")
}

func dot(w *world, _ []map[string]object) error {
	if len(w.stack) == 0 {
		return emptyStack
	}
	n, ok := w.pop().(numo)
	if !ok {
		return errors.New("Only numbers can be dot-printed.")
	}
	_, err := fmt.Fprintf(w.output, "%d\n", int(n))
	return err
}

func plus(w *world, _ []map[string]object) error {
	if len(w.stack) < 2 {
		return emptyStack
	}
	n1, ok := w.pop().(numo)
	if !ok {
		return errors.New("Only numbers can be added.")
	}
	n2, ok := w.pop().(numo)
	if !ok {
		return errors.New("Only numbers can be added.")
	}
	w.push(n1 + n2)
	return nil
}

func negate(w *world, _ []map[string]object) error {
	if len(w.stack) == 0 {
		return emptyStack
	}
	n, ok := w.pop().(numo)
	if !ok {
		return errors.New("Only numbers can be negated.")
	}
	w.push(-n)
	return nil
}

func isZero(w *world, _ []map[string]object) error {
	if len(w.stack) == 0 {
		return emptyStack
	}
	n, ok := w.pop().(numo)
	if !ok {
		return errors.New("Only numbers can be checked for zeroness.")
	}
	if n == 0 {
		w.push(innerTrue)
	} else {
		w.push(innerFalse)
	}
	return nil
}

func mod(w *world, _ []map[string]object) error {
	if len(w.stack) < 2 {
		return emptyStack
	}
	n1, ok := w.pop().(numo)
	if !ok {
		return errors.New("Mod is for numbers.")
	}
	n2, ok := w.pop().(numo)
	if !ok {
		return errors.New("Mod is for numbers.")
	}
	w.push(n2 % n1)
	return nil
}

func cDot(w *world, _ []map[string]object) error {
	if len(w.stack) == 0 {
		return emptyStack
	}
	c, ok := w.pop().(charo)
	if !ok {
		return errors.New("Only chars can be char-printed.")
	}
	_, err := fmt.Fprintf(w.output, "%c", rune(c))
	return err
}

func cEq(w *world, _ []map[string]object) error {
	if len(w.stack) < 2 {
		return emptyStack
	}
	c1, ok := w.pop().(charo)
	if !ok {
		return errors.New("Only chars can be char-compared.")
	}
	c2, ok := w.pop().(charo)
	if !ok {
		return errors.New("Only chars can be char-compared.")
	}
	if (c1 == c2) {
		w.push(innerTrue)
	} else {
		w.push(innerFalse)
	}
	return nil
}

func sIndex(w *world, _ []map[string]object) error {
	if len(w.stack) < 2 {
		return emptyStack
	}
	i, ok := w.pop().(numo)
	if !ok {
		return errors.New("String index must be a number.")
	}
	s, ok := w.pop().(stringo)
	if !ok {
		return errors.New("Only strings can be string-indexed.")
	}
	w.push(charo(s[i]))
	return nil
}

func sLen(w *world, _ []map[string]object) error {
	if len(w.stack) == 0 {
		return emptyStack
	}
	s, ok := w.pop().(stringo)
	if !ok {
		return errors.New("Only strings have string length.")
	}
	w.push(numo(len(s)))
	return nil
}

var innerTrue = builtin(func(w *world, context []map[string]object) error {
	if len(w.stack) < 2 {
		return emptyStack
	}
	w.pop()
	return bangc.eval(bangc{}, w, context)
})

var innerFalse = builtin(func(w *world, context []map[string]object) error {
	if len(w.stack) < 2 {
		return emptyStack
	}
	f := w.pop()
	w.pop()
	w.push(f)
	return bangc.eval(bangc{}, w, context)
})
