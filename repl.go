package main

import (
	"bufio"
	"fmt"
	. "github.com/jonnyarnold/fn-go/parser"
	. "github.com/jonnyarnold/fn-go/runtime"
	. "github.com/jonnyarnold/fn-go/tokeniser"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		text, _ := reader.ReadString('\n')

		tokens := Tokenise(text)
		expressions, _ := Parse(tokens)
		result := Execute(expressions)

		if result.Error != nil {
			fmt.Println(result.Error.Error())
		}

		if result.Value != nil {
			fmt.Println(result.Value.String())
		}
	}
}
