package vm

import (
	"strconv"
)

type vmBool struct{ Value bool }

func (b vmBool) String() string {
	return strconv.FormatBool(b.Value)
}

func (b vmBool) Negate() vmConstant {
	return vmBool{Value: !b.Value}
}

func (b vmBool) IsFalse() bool {
	return !b.Value
}
