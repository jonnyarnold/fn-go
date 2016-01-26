package vm

// constantValue is a structure (because there's no union in Go)
// that holds one of the primitive types.
type constantValue struct {
	Integer byte
	// Float   float64
	// String  string
	Boolean bool
}

const (
	TYPE_INT    = 0
	TYPE_FLOAT  = 1
	TYPE_STRING = 2
	TYPE_BOOL   = 3
)

// Converts a []byte to a constantValue
func ConstantValue(bytes []byte) constantValue {
	switch bytes[0] {
	case TYPE_INT:
		return constantValue{Integer: bytes[1]}
	case TYPE_BOOL:
		return constantValue{Boolean: bytes[1] != 0}
	}

	panic("Could not constantize!")
}
