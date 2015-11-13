package parser

type Expression interface{}

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

// A function prototype (a literal).
type FunctionPrototypeExpression struct {
	Arguments []IdentifierExpression
	Body      BlockExpression
}

// A function call.
type FunctionCallExpression struct {
	Identifier IdentifierExpression
	Arguments  []Expression
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
