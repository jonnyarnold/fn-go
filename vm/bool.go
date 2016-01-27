package vm

import (
	"strconv"
)

type vmBool struct{ Value bool }

func (b vmBool) String() string {
	return strconv.FormatBool(b.Value)
}
