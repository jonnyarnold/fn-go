package repl

import (
	"bufio"
	"fmt"
	. "github.com/jonnyarnold/fn-go/compiler/parser"
	. "github.com/jonnyarnold/fn-go/compiler/runtime"
	. "github.com/jonnyarnold/fn-go/compiler/tokeniser"
	"os"
)

// Starts the REPL and takes control.
func Run() {
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
