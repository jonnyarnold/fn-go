package compiler

import (
	"fmt"
	. "github.com/jonnyarnold/fn-go/compiler/parser"
	// . "github.com/jonnyarnold/fn-go/compiler/runtime"
	. "github.com/jonnyarnold/fn-go/compiler/tokeniser"
	. "github.com/jonnyarnold/fn-go/compiler/transpiler"
	. "github.com/jonnyarnold/fn-go/vm"
	"io/ioutil"
)

func Run(fileName string) {
	file, _ := ioutil.ReadFile(fileName)

	tokens := Tokenise(string(file))

	// for _, token := range tokens {
	// 	fmt.Println(token)
	// }

	expressions, errors := Parse(tokens)

	if errors != nil {
		fmt.Println(errors)
	}

	for _, expr := range expressions {
		fmt.Println(expr)
	}

	result := Transpile(expressions)
	for _, b := range result {
		fmt.Printf("%02x ", b)
	}
	fmt.Print("\n")

	machine := NewMachine(result)
	machine.ProcessAll()
	fmt.Println(machine.String())
	fmt.Println(machine.Result())

	// result := Execute(expressions)

	// if result.Error != nil {
	// 	fmt.Println(result.Error)
	// }
}
