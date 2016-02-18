package vm

type vmConstant interface {
	String() string
	Negate() vmConstant
	IsFalse() bool
}

const (
	TYPE_INT    = 0
	TYPE_FLOAT  = 1
	TYPE_STRING = 2
	TYPE_BOOL   = 3
)

// Converts a []byte to a vmConstant
func VMConstant(bytes []byte) vmConstant {
	switch bytes[0] {
	case TYPE_INT:
		return vmNumber{Type: TYPE_INT, Integer: int(bytes[1])}
	case TYPE_FLOAT:
		return vmNumber{Type: TYPE_FLOAT, Float: float64(bytes[1])}
	case TYPE_BOOL:
		return vmBool{Value: bytes[1] != 0}
	}

	panic("Could not constantize!")
}
