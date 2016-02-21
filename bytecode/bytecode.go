package bytecode

import (
	"bytes"
	"encoding/binary"
	// "fmt"
)

type Bytecode []byte

func NewBytecode() Bytecode {
	return Bytecode{}
}

func (b Bytecode) DeclareBool(value bool) Bytecode {
	var boolType byte
	if value {
		boolType = TYPE_TRUE
	} else {
		boolType = TYPE_FALSE
	}

	return append(b, []byte{DECLARE, boolType}...)
}

func (b Bytecode) DeclareInt(value int64) Bytecode {
	valueBytes := make([]byte, 8)
	binary.PutVarint(valueBytes, value)

	bytes := append(b, []byte{DECLARE, TYPE_INT}...)
	bytes = append(bytes, valueBytes[0:8]...)

	return bytes
}

func (b Bytecode) DeclareFloat(value float64) Bytecode {
	valueBytes := bytes.NewBuffer([]byte{})
	binary.Write(valueBytes, binary.LittleEndian, value)

	bytes := append(b, []byte{DECLARE, TYPE_FLOAT}...)
	bytes = append(bytes, valueBytes.Bytes()...)

	return bytes
}

func (b Bytecode) DeclareString(value string) Bytecode {
	valueLength := len(value)
	valueLengthBytes := make([]byte, 8)
	binary.PutUvarint(valueLengthBytes, uint64(valueLength))

	bytes := append(b, []byte{DECLARE, TYPE_STRING}...)
	bytes = append(bytes, valueLengthBytes...)
	bytes = append(bytes, []byte(value)...)

	return bytes
}

func (b Bytecode) Copy(constantIndex byte) Bytecode {
	return append(b, []byte{COPY, constantIndex}...)
}

func (b Bytecode) Add() Bytecode {
	return append(b, ADD)
}

func (b Bytecode) Subtract() Bytecode {
	return append(b, SUBTRACT)
}

func (b Bytecode) Multiply() Bytecode {
	return append(b, MULTIPLY)
}

func (b Bytecode) Divide() Bytecode {
	return append(b, DIVIDE)
}

func (b Bytecode) And() Bytecode {
	return append(b, AND)
}

func (b Bytecode) Or() Bytecode {
	return append(b, OR)
}

func (b Bytecode) Negate() Bytecode {
	return append(b, NEGATE)
}

func (b Bytecode) JumpIfFalse(innerCode Bytecode) Bytecode {
	jumpCount := len(innerCode)
	jumpCountBytes := make([]byte, 8)
	binary.PutUvarint(jumpCountBytes, uint64(jumpCount))

	bytes := append(b, JUMP_IF_FALSE)
	bytes = append(bytes, jumpCountBytes...)
	bytes = append(bytes, innerCode...)
	return bytes
}

func (b Bytecode) FunctionFollows(innerCode Bytecode) Bytecode {
	jumpCount := len(innerCode)
	jumpCountBytes := make([]byte, 8)
	binary.PutUvarint(jumpCountBytes, uint64(jumpCount))

	bytes := append(b, FUNCTION_FOLLOWS)
	bytes = append(bytes, jumpCountBytes...)
	bytes = append(bytes, innerCode...)
	return bytes
}

func (b Bytecode) Call(constIdx byte) Bytecode {
	return append(b, CALL, constIdx)
}

func (b Bytecode) Return() Bytecode {
	return append(b, RETURN)
}

func (b Bytecode) Dump() Bytecode {
	return append(b, DUMP)
}
