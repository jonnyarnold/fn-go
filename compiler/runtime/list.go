package runtime

import (
	"bytes"
	"errors"
	"fmt"
)

// Alias for a function that can be used within fn.
type list struct {
	Items []fnScope
}

func (list list) Definitions() defMap {
	allDefs := defMap{
		"each":     fn([]string{"fn"}, list.each),
		"asString": fn([]string{}, list.asString),
	}

	for key, value := range DefaultScope().Definitions() {

		_, ok := allDefs[key]
		if !ok {
			allDefs[key] = value
		}
	}

	return allDefs
}

func (list list) Define(id string, value fnScope) (fnScope, error) {
	panic("Attempted defining " + id + " on a list!")
}

func (list list) String() string {
	var str bytes.Buffer
	str.WriteString("List(")

	lastIdx := len(list.Items) - 1
	for idx, item := range list.Items {
		str.WriteString(item.String())
		if idx != lastIdx {
			str.WriteString(", ")
		}
	}

	str.WriteString(")")
	return str.String()
}

func (list list) Call(args []fnScope) (fnScope, error) {
	if len(args) != 1 {
		return nil, errors.New(fmt.Sprintf(
			"Argument number mismatch: got %i, need 1",
			len(args),
		))
	}

	index := args[0].(number).AsInt()
	return list.Items[index], nil
}

func (list list) Value() interface{} {
	return list.Items
}

func List(values []fnScope) (fnScope, error) {
	return list{Items: values}, nil
}

func (list list) each(args []fnScope) (fnScope, error) {
	for _, item := range list.Items {
		_, err := args[0].Call([]fnScope{item})
		if err != nil {
			return nil, err
		}
	}

	return nil, nil
}

func (list list) asString(args []fnScope) (fnScope, error) {
	return FnString(list.String()), nil
}

type fnList struct{}

func (list fnList) Definitions() defMap {
	return defaultScope{}.Definitions()
}

func (list fnList) Define(id string, value fnScope) (fnScope, error) {
	panic("Attempted definition on the List function!")
}

func (list fnList) String() string {
	return "List(...)"
}

func (list fnList) Call(args []fnScope) (fnScope, error) {
	return List(args)
}

func (list fnList) Value() interface{} {
	return nil
}
