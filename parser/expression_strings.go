package parser

import (
	"bytes"
	"fmt"
)

func (ne NumberExpression) String() string {
	return fmt.Sprint(ne.Value)
}

func (se StringExpression) String() string {
	return fmt.Sprint(se.Value)
}

func (be BooleanExpression) String() string {
	return fmt.Sprintf("%t", be.Value)
}

func (ie IdentifierExpression) String() string {
	return fmt.Sprintf(ie.Name)
}

func (be BlockExpression) String() string {
	var str bytes.Buffer
	str.WriteString("{\n")

	for _, expr := range be.Body {
		str.WriteString("  ")
		str.WriteString(expr.String())
		str.WriteString("\n")
	}

	str.WriteString("}")
	return str.String()
}

func (args arguments) String() string {
	var str bytes.Buffer
	str.WriteString("(")

	lastIdx := len(args) - 1
	for idx, arg := range args {
		str.WriteString(arg.String())
		if idx != lastIdx {
			str.WriteString(", ")
		}
	}

	str.WriteString(")")
	return str.String()
}

func (fpe FunctionPrototypeExpression) String() string {
	return fmt.Sprintf(
		"%s %s",
		fpe.Arguments.String(),
		fpe.Body.String(),
	)
}

func (params params) String() string {
	var str bytes.Buffer
	str.WriteString("(")

	lastIdx := len(params) - 1
	for idx, param := range params {
		str.WriteString(param.String())
		if idx != lastIdx {
			str.WriteString(", ")
		}
	}

	str.WriteString(")")
	return str.String()
}

func (fce FunctionCallExpression) String() string {
	return fmt.Sprintf(
		"%s%s",
		fce.Identifier.String(),
		fce.Arguments.String(),
	)
}

func (ce ConditionalExpression) String() string {
	var str bytes.Buffer
	str.WriteString("when {\n")

	for _, cond := range ce.Branches {
		str.WriteString("  ")
		str.WriteString(cond.String())
		str.WriteString("\n")
	}

	str.WriteString("}")
	return str.String()
}

func (cbe ConditionalBranchExpression) String() string {
	return fmt.Sprintf(
		"%s %s",
		cbe.Condition.String(),
		cbe.Body.String(),
	)
}
