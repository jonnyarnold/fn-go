package transpiler

import (
	"fmt"
	. "github.com/jonnyarnold/fn-go/bytecode"
	. "github.com/jonnyarnold/fn-go/compiler/parser"
	"strconv"
	"strings"
)

// Converts the Abstract Syntax Tree
// into bytecode that can be run by the VM.
func Transpile(ast []Expression) []byte {
	t := transpiler{ast: ast, output: NewBytecode()}
	t.transpileAst()
	return t.output
}

type transpiler struct {
	ast            []Expression
	symbols        []string
	blockNameStack []string
	output         Bytecode
}

func (t *transpiler) transpileAst() {
	for _, expression := range t.ast {
		t.output = append(t.output, t.transpile(expression)...)
	}

	t.output = t.output.Return()
}

func (t *transpiler) transpile(e Expression) Bytecode {
	switch e.(type) {

	case BooleanExpression:
		value := e.(BooleanExpression).Value
		return NewBytecode().DeclareBool(value)

	case StringExpression:
		value := e.(StringExpression).Value
		return NewBytecode().DeclareString(value)

	case NumberExpression:
		// We might have ourselves a float or an integer
		strValue := e.(NumberExpression).Value

		if strings.ContainsRune(strValue, '.') {
			floatValue, _ := strconv.ParseFloat(strValue, 64)
			return NewBytecode().DeclareFloat(floatValue)
		} else {
			intValue, _ := strconv.ParseInt(strValue, 10, 64)
			return NewBytecode().DeclareInt(intValue)
		}

	case IdentifierExpression:
		name := e.(IdentifierExpression).Name

		var constIdx = -1
		for idx, symbol := range t.symbols {
			if name == symbol {
				constIdx = idx
				break
			}
		}

		if constIdx == -1 {
			panic(fmt.Sprintf("Cound not resolve identifier %v", name))
		}

		return NewBytecode().Copy(byte(constIdx))

	case FunctionCallExpression:
		fnName := e.(FunctionCallExpression).Identifier.Name
		switch fnName {
		case "=":
			id := e.(FunctionCallExpression).Arguments[0].(IdentifierExpression).Name
			value := e.(FunctionCallExpression).Arguments[1]

			t.symbols = append(t.symbols, id)
			return t.transpile(value)

		case "+":
			first := e.(FunctionCallExpression).Arguments[0]
			second := e.(FunctionCallExpression).Arguments[1]

			code := append(t.transpile(first), t.transpile(second)...)
			code = append(code, NewBytecode().Add()...)
			return code

		default:
			panic("Unknown function name.")
		}

	default:
		panic("Cannot convert expression to bytecode.")
	}
}
