package eval

import (
	"testing"
)

func TestEvalInterface(t *testing.T) {
	// Not an actual test, just check that we are satisfying the interface
	dummy := func(e Evaluator) {}
	dummy(Eval)
}
