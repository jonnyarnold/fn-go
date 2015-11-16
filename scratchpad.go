package main

import (
	"fmt"
	. "github.com/jonnyarnold/fn-go/parser"
	. "github.com/jonnyarnold/fn-go/tokeniser"
	"io/ioutil"
)

func main() {
	file, _ := ioutil.ReadFile("tour.fn")

	tokens := Tokenise(string(file))

	// for _, token := range tokens {
	// 	fmt.Println(token)
	// }

	expressions, errors := Parse(tokens)

	for _, expression := range expressions {
		fmt.Println(expression)
	}

	if errors != nil {
		fmt.Println(errors)
	}
}
