package runtime

import (
	. "github.com/jonnyarnold/fn-go/parser"
)

func execNumber(expr NumberExpression) number {
	return Number(expr.Value)
}

func execString(expr StringExpression) fnString {
	return FnString(expr.Value)
}

func execBool(expr BooleanExpression) fnBool {
	return FnBool(expr.Value)
}
