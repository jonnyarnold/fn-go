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
	return defMap{}
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
