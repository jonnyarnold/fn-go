package main

import (
	"bufio"
	"fmt"
	. "github.com/jonnyarnold/fn-go/parser"
	. "github.com/jonnyarnold/fn-go/runtime"
	. "github.com/jonnyarnold/fn-go/tokeniser"
	"os"
)

func repl() {
	reader := bufio.NewReader(os.Stdin)
	replScope := DefaultScope()

	for {
		fmt.Print("> ")
		text, _ := reader.ReadString('\n')

		tokens := Tokenise(text)
		expressions, _ := Parse(tokens)
		result := ExecuteIn(expressions, replScope)

		if result.Error != nil {
			fmt.Println(result.Error.Error())
		}

		if result.Value != nil {
			fmt.Println(result.Value.String())
		}
	}
}
