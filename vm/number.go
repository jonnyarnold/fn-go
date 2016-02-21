package vm

import (
	"fmt"
	. "github.com/jonnyarnold/fn-go/bytecode"
)

type vmNumber struct {
	Type    uint
	Integer int64
	Float   float64
}

func (n vmNumber) String() string {
	switch n.Type {
	case TYPE_INT:
		return fmt.Sprintf("%v", n.Integer)
	case TYPE_FLOAT:
		return fmt.Sprintf("%v", n.Float)
	}

	panic("Unknown type!")
}

func (n vmNumber) AsFloat() float64 {
	switch n.Type {
	case TYPE_INT:
		return float64(n.Integer)
	case TYPE_FLOAT:
		return n.Float
	}

	panic("Unknown type!")
}

func (n vmNumber) AsInt() int64 {
	switch n.Type {
	case TYPE_INT:
		return n.Integer
	case TYPE_FLOAT:
		return int64(n.Float)
	}

	panic("Unknown type!")
}

func (n vmNumber) IsFalse() bool {
	return false
}

func AddNumbers(first vmNumber, second vmNumber) vmNumber {
	if first.Type == second.Type && first.Type == TYPE_INT {
		return vmNumber{Type: TYPE_INT, Integer: first.Integer + second.Integer}
	}

	if first.Type == second.Type && first.Type == TYPE_FLOAT {
		return vmNumber{Type: TYPE_FLOAT, Float: first.Float + second.Float}
	}

	return vmNumber{Type: TYPE_FLOAT, Float: first.AsFloat() + second.AsFloat()}
}

func SubtractNumbers(first vmNumber, second vmNumber) vmNumber {
	if first.Type == second.Type && first.Type == TYPE_INT {
		return vmNumber{Type: TYPE_INT, Integer: first.Integer - second.Integer}
	}

	if first.Type == second.Type && first.Type == TYPE_FLOAT {
		return vmNumber{Type: TYPE_FLOAT, Float: first.Float - second.Float}
	}

	return vmNumber{Type: TYPE_FLOAT, Float: first.AsFloat() - second.AsFloat()}
}

func MultiplyNumbers(first vmNumber, second vmNumber) vmNumber {
	if first.Type == second.Type && first.Type == TYPE_INT {
		return vmNumber{Type: TYPE_INT, Integer: first.Integer * second.Integer}
	}

	if first.Type == second.Type && first.Type == TYPE_FLOAT {
		return vmNumber{Type: TYPE_FLOAT, Float: first.Float * second.Float}
	}

	return vmNumber{Type: TYPE_FLOAT, Float: first.AsFloat() * second.AsFloat()}
}

func DivideNumbers(first vmNumber, second vmNumber) vmNumber {
	return vmNumber{Type: TYPE_FLOAT, Float: first.AsFloat() / second.AsFloat()}
}

func (n vmNumber) Negate() vmConstant {
	return vmNumber{Type: n.Type, Integer: -n.Integer, Float: -n.Float}
}
