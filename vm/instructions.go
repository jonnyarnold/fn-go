package vm

// Constants to allow reference to the instructions.
const (
	// Debugging
	DUMP = iota

	// Declarations
	DECLARE          = iota
	FUNCTION_FOLLOWS = iota
	COPY             = iota

	// Artihmetic operations
	ADD      = iota
	SUBTRACT = iota
	MULTIPLY = iota
	DIVIDE   = iota

	// Boolean operations
	AND = iota
	OR  = iota

	// Unary operations
	NEGATE = iota

	// Control flow
	JUMP_IF_FALSE = iota
	CALL          = iota
	RETURN        = iota
)
