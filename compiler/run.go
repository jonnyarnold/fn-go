package compiler

import (
	"fmt"
	. "github.com/jonnyarnold/fn-go/compiler/parser"
	. "github.com/jonnyarnold/fn-go/compiler/runtime"
	. "github.com/jonnyarnold/fn-go/compiler/tokeniser"
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

	// for _, expr := range expressions {
	// 	fmt.Println(expr)
	// }

	result := Execute(expressions)

	if result.Error != nil {
		fmt.Println(result.Error)
	}
}
