package tokeniser

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func SoCodeYieldsTokens(code string, tokens []Token) {
	tokenised := Tokenise(code)

	// Lines and columns are output.
	// We don't care.
	for idx, _ := range tokenised {
		tokenised[idx].Line = 0
		tokenised[idx].Column = 0
	}

	So(tokenised, ShouldResemble, tokens)
}

func TestTokeniser(t *testing.T) {
	Convey("Spaces", t, func() {
		Convey("are ignored", func() {
			SoCodeYieldsTokens("  \n", []Token{})
		})
	})

	Convey("Comments", t, func() {
		Convey("are ignored", func() {
			SoCodeYieldsTokens("# This is a comment!", []Token{})
			SoCodeYieldsTokens("# This is a comment!\n", []Token{})
		})

		Convey("are only ignored after the hash symbol", func() {
			SoCodeYieldsTokens("a # This is a comment!\n", []Token{
				Token{Type: "identifier", Value: "a"},
			})
		})

		Convey("stop at a newline", func() {
			SoCodeYieldsTokens("# This is a comment!\nb", []Token{
				Token{Type: "identifier", Value: "b"},
			})
		})
	})

	Convey("Open bracket is found", t, func() {
		SoCodeYieldsTokens("(", []Token{
			Token{Type: "bracket_open"},
		})
	})

	Convey("Close bracket is found", t, func() {
		SoCodeYieldsTokens(")", []Token{
			Token{Type: "bracket_close"},
		})
	})

	Convey("Comma is found", t, func() {
		SoCodeYieldsTokens(",", []Token{
			Token{Type: "comma"},
		})
	})

	Convey("End statement is found", t, func() {
		SoCodeYieldsTokens(";", []Token{
			Token{Type: "end_statement"},
		})
	})

	Convey("Use is its own token", t, func() {
		SoCodeYieldsTokens("us use user", []Token{
			Token{Type: "identifier", Value: "us"},
			Token{Type: "use"},
			Token{Type: "identifier", Value: "user"},
		})
	})

	Convey("Import is its own token", t, func() {
		SoCodeYieldsTokens("impor import imports", []Token{
			Token{Type: "identifier", Value: "impor"},
			Token{Type: "import"},
			Token{Type: "identifier", Value: "imports"},
		})
	})

	Convey("When is its own token", t, func() {
		SoCodeYieldsTokens("whe when whent", []Token{
			Token{Type: "identifier", Value: "whe"},
			Token{Type: "when"},
			Token{Type: "identifier", Value: "whent"},
		})
	})

	Convey("Infix operators", t, func() {
		Convey("are found with spaces around them", func() {
			SoCodeYieldsTokens("a eq b", []Token{
				Token{Type: "identifier", Value: "a"},
				Token{Type: "infix_operator", Value: "eq"},
				Token{Type: "identifier", Value: "b"},
			})
		})

		Convey("without symbols are not found inside identifiers", func() {
			SoCodeYieldsTokens("aeqb", []Token{
				Token{Type: "identifier", Value: "aeqb"},
			})
		})

		Convey("with symbols are found inside identifiers", func() {
			SoCodeYieldsTokens("a/b", []Token{
				Token{Type: "identifier", Value: "a"},
				Token{Type: "infix_operator", Value: "/"},
				Token{Type: "identifier", Value: "b"},
			})
		})
	})

	Convey("Open block is found", t, func() {
		SoCodeYieldsTokens("{", []Token{
			Token{Type: "block_open"},
		})
	})

	Convey("Close block is found", t, func() {
		SoCodeYieldsTokens("}", []Token{
			Token{Type: "block_close"},
		})
	})

	Convey("Strings", t, func() {
		Convey("are found with double quotes", func() {
			SoCodeYieldsTokens("\"Hello!\"", []Token{
				Token{Type: "string", Value: "Hello!"},
			})
		})

		Convey("can escape double quotes within them", nil)
	})

	Convey("Numbers", t, func() {

		Convey("can be integers", func() {
			SoCodeYieldsTokens("1", []Token{
				Token{Type: "number", Value: "1"},
			})
		})

		Convey("can be floats", func() {
			SoCodeYieldsTokens("1.5", []Token{
				Token{Type: "number", Value: "1.5"},
			})
		})

	})

	Convey("Booleans", t, func() {

		Convey("include true", func() {
			SoCodeYieldsTokens("tru true truer", []Token{
				Token{Type: "identifier", Value: "tru"},
				Token{Type: "boolean", Value: "true"},
				Token{Type: "identifier", Value: "truer"},
			})
		})

		Convey("include false", func() {
			SoCodeYieldsTokens("fals false falser", []Token{
				Token{Type: "identifier", Value: "fals"},
				Token{Type: "boolean", Value: "false"},
				Token{Type: "identifier", Value: "falser"},
			})
		})

	})

	Convey("Identifiers", t, func() {

		Convey("hashes", func() {

			Convey("can not start with a hash", func() {
				SoCodeYieldsTokens("#id", []Token{})
			})

			Convey("can not include a hash", func() {
				SoCodeYieldsTokens("ab#cd", []Token{
					Token{Type: "identifier", Value: "ab"},
				})
			})

		})

		Convey("brackets", func() {

			Convey("can not include an opening bracket", func() {
				SoCodeYieldsTokens("ab(cd", []Token{
					Token{Type: "identifier", Value: "ab"},
					Token{Type: "bracket_open"},
					Token{Type: "identifier", Value: "cd"},
				})
			})

			Convey("can not include a closing bracket", func() {
				SoCodeYieldsTokens("ab)cd", []Token{
					Token{Type: "identifier", Value: "ab"},
					Token{Type: "bracket_close"},
					Token{Type: "identifier", Value: "cd"},
				})
			})
		})

		Convey("can not include a comma", func() {
			SoCodeYieldsTokens("ab,cd", []Token{
				Token{Type: "identifier", Value: "ab"},
				Token{Type: "comma"},
				Token{Type: "identifier", Value: "cd"},
			})
		})

		Convey("can not include an end statement", func() {
			SoCodeYieldsTokens("ab;cd", []Token{
				Token{Type: "identifier", Value: "ab"},
				Token{Type: "end_statement"},
				Token{Type: "identifier", Value: "cd"},
			})
		})

		Convey("blocks", func() {
			Convey("can not include an opening block", func() {
				SoCodeYieldsTokens("ab{cd", []Token{
					Token{Type: "identifier", Value: "ab"},
					Token{Type: "block_open"},
					Token{Type: "identifier", Value: "cd"},
				})
			})

			Convey("can not include a closing block", func() {
				SoCodeYieldsTokens("ab}cd", []Token{
					Token{Type: "identifier", Value: "ab"},
					Token{Type: "block_close"},
					Token{Type: "identifier", Value: "cd"},
				})
			})
		})

		Convey("can not include a double quote", func() {
			SoCodeYieldsTokens("ab\"cd", []Token{
				Token{Type: "identifier", Value: "ab"},
				Token{Type: "string", Value: "cd"},
			})
		})

		Convey("can not include a space", func() {
			SoCodeYieldsTokens("ab cd", []Token{
				Token{Type: "identifier", Value: "ab"},
				Token{Type: "identifier", Value: "cd"},
			})
		})

		Convey("numbers", func() {
			Convey("can not start with a number", func() {
				SoCodeYieldsTokens("1a", []Token{
					Token{Type: "number", Value: "1"},
					Token{Type: "identifier", Value: "a"},
				})
			})

			Convey("can include a number", func() {
				SoCodeYieldsTokens("a1", []Token{
					Token{Type: "identifier", Value: "a1"},
				})
			})
		})

		Convey("symbols", func() {
			Convey("can start with symbols", func() {
				SoCodeYieldsTokens("@", []Token{
					Token{Type: "identifier", Value: "@"},
				})
			})

			Convey("can include symbols", func() {
				SoCodeYieldsTokens("a@b", []Token{
					Token{Type: "identifier", Value: "a@b"},
				})
			})
		})
	})
}
