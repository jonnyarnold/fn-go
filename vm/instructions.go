package vm

// Set of instruction types recognised by the VM.
const (
	TERMINATE                  = 1
	DECLARE_CONSTANT           = 2
	SET_REGISTER_WITH_CONSTANT = 3
	ADD                        = 4
	SUBTRACT                   = 5
	MULTIPLY                   = 6
	DIVIDE                     = 7
	DUMP                       = 0
)
