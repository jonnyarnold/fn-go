package vm

type vmString struct{ Value string }

func (s vmString) String() string {
	return s.Value
}

func (s vmString) Negate() vmConstant {
	panic("Attempted NEGATE on a string!")
}

func (s vmString) IsFalse() bool {
	return false
}
