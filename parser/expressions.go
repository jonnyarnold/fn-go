package parser

type Expression interface {
	String() string
}

// A number literal.
type NumberExpression struct {
	Value string // Kept as a String, as it may be integer or float.
}

// A string literal.
type StringExpression struct {
	Value string
}

// A boolean literal.
type BooleanExpression struct {
	Value bool
}

// An identifier.
type IdentifierExpression struct {
	Name string
}

// A block expression (a literal).
type BlockExpression struct {
	Body []Expression
}

type arguments []IdentifierExpression

// A function prototype (a literal).
type FunctionPrototypeExpression struct {
	Arguments arguments
	Body      BlockExpression
}

type params []Expression

// A function call.
type FunctionCallExpression struct {
	Identifier IdentifierExpression
	Arguments  params
}

// A conditional expression.
type ConditionalExpression struct {
	Branches []ConditionalBranchExpression
}

// A branch of a conditional expression.
type ConditionalBranchExpression struct {
	Condition Expression
	Body      BlockExpression
}
