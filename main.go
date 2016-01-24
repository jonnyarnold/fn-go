package main

import (
	"fmt"
	. "github.com/jonnyarnold/fn-go/parser"
	. "github.com/jonnyarnold/fn-go/runtime"
	. "github.com/jonnyarnold/fn-go/tokeniser"
	"io/ioutil"
	"os"
)

func main() {
	file, _ := ioutil.ReadFile(os.Args[1])

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
