package runtime

import (
	. "github.com/jonnyarnold/fn-go/parser"
)

// Executes the given expressions in the default scope.
func Execute(exprs []Expression) error {
	_, err := ExecuteIn(exprs, DefaultScope)
	return err
}
