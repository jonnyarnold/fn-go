package vm

import (
	"bytes"
	"encoding/binary"
	. "github.com/jonnyarnold/fn-go/bytecode"
	"math"
)

type vmConstant interface {
	String() string
	Negate() vmConstant
	IsFalse() bool
}

// Converts a []byte to a vmConstant
func VMConstant(constType byte, valueBytes []byte) vmConstant {
	switch constType {
	case TYPE_INT:
		byteBuffer := bytes.NewReader(valueBytes)
		value, err := binary.ReadVarint(byteBuffer)
		if err != nil {
			panic(err)
		}

		return vmNumber{Type: TYPE_INT, Integer: value}
	case TYPE_FLOAT:
		valueBits := binary.LittleEndian.Uint64(valueBytes)
		value := math.Float64frombits(valueBits)

		return vmNumber{Type: TYPE_FLOAT, Float: value}
	case TYPE_TRUE:
		return vmBool{Value: true}
	case TYPE_FALSE:
		return vmBool{Value: false}
	case TYPE_STRING:
		return vmString{Value: string(valueBytes)}
	}

	panic("Could not constantize!")
}
