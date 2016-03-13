package runtime

import (
	"errors"
	"strconv"
)

// A number is a scope representing a numeric value.
type number struct {
	value string
}

func (num number) Definitions() defMap {
	return defMap{
		"+":        fn([]string{"other"}, num.add),
		"-":        fn([]string{"other"}, num.subtract),
		"*":        fn([]string{"other"}, num.multiply),
		"/":        fn([]string{"other"}, num.divide),
		"and":      fn([]string{"other"}, num.and),
		"or":       fn([]string{"other"}, num.or),
		"eq":       fn([]string{"other"}, num.eq),
		"moreThan": fn([]string{"other"}, num.moreThan),
		"lessThan": fn([]string{"other"}, num.lessThan),
	}
}

func (num number) Define(id string, value fnScope) (fnScope, error) {
	return nil, errors.New("Attempted definition on a number!")
}

func (num number) String() string {
	return num.value
}

func (num number) Call(args []fnScope) (fnScope, error) {
	return nil, errors.New("Number called as a function!")
}

func (num number) Value() interface{} {
	return num.AsFloat()
}

func (num number) AsFloat() float64 {
	f, err := strconv.ParseFloat(num.value, 64)
	if err != nil {
		panic(err)
	}

	return f
}

func (num number) AsInt() int64 {
	i, err := strconv.ParseInt(num.value, 10, 64)
	if err != nil {
		panic(err)
	}

	return i
}

func Number(num string) number {
	return number{value: num}
}

func NumberFromFloat(num float64) number {
	return Number(strconv.FormatFloat(num, 'f', -1, 64))
}

func (num number) add(args []fnScope) (fnScope, error) {
	return NumberFromFloat(num.AsFloat() + args[0].(number).AsFloat()), nil
}

func (num number) subtract(args []fnScope) (fnScope, error) {
	return NumberFromFloat(num.AsFloat() - args[0].(number).AsFloat()), nil
}

func (num number) multiply(args []fnScope) (fnScope, error) {
	return NumberFromFloat(num.AsFloat() * args[0].(number).AsFloat()), nil
}

func (num number) divide(args []fnScope) (fnScope, error) {
	return NumberFromFloat(num.AsFloat() / args[0].(number).AsFloat()), nil
}

func (num number) moreThan(args []fnScope) (fnScope, error) {
	return FnBool(num.AsFloat() > args[0].(number).AsFloat()), nil
}

func (num number) lessThan(args []fnScope) (fnScope, error) {
	return FnBool(num.AsFloat() < args[0].(number).AsFloat()), nil
}

func (self number) and(args []fnScope) (fnScope, error) {
	return FnBool(AsBool(self) && AsBool(args[0])), nil
}

func (self number) or(args []fnScope) (fnScope, error) {
	return FnBool(AsBool(self) || AsBool(args[0])), nil
}

func (self number) eq(args []fnScope) (fnScope, error) {
	return FnBool(self.Value() == args[0].Value()), nil
}
