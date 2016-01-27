package vm

// Constants to allow reference to the instructions.
const (
	// Program control
	TERMINATE = iota
	DUMP      = iota

	// Declarations
	DECLARE_CONSTANT           = iota
	SET_REGISTER_WITH_CONSTANT = iota

	// Artihmetic operations
	ADD      = iota
	SUBTRACT = iota
	MULTIPLY = iota
	DIVIDE   = iota

	// Boolean operations
	AND = iota
	OR  = iota
	NOT = iota
)
