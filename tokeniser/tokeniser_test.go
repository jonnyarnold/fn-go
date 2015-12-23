package tokeniser

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestTokeniser(t *testing.T) {
	Convey("Spaces", t, func() {
		Convey("are ignored", func() {
			So(Tokenise("  \n"), ShouldResemble, []Token{})
		})
	})

	Convey("Comments", t, func() {
		Convey("are ignored", func() {
			So(Tokenise("# This is a comment!\n"), ShouldResemble, []Token{})
		})

		Convey("are only ignored after the hash symbol", func() {
			So(Tokenise("a # This is a comment!\n"), ShouldResemble, []Token{
				Token{Type: "identifier", Value: "a"},
			})
		})

		Convey("stop at a newline", func() {
			So(Tokenise("# This is a comment!\nb"), ShouldResemble, []Token{
				Token{Type: "identifier", Value: "b"},
			})
		})
	})

	Convey("Open bracket is found", t, func() {
		So(Tokenise("("), ShouldResemble, []Token{
			Token{Type: "bracket_open"},
		})
	})

	Convey("Close bracket is found", t, func() {
		So(Tokenise(")"), ShouldResemble, []Token{
			Token{Type: "bracket_close"},
		})
	})

	Convey("Comma is found", t, func() {
		So(Tokenise(","), ShouldResemble, []Token{
			Token{Type: "comma"},
		})
	})

	Convey("End statement is found", t, func() {
		So(Tokenise(";"), ShouldResemble, []Token{
			Token{Type: "end_statement"},
		})
	})

	Convey("Use is its own token", t, func() {
		So(Tokenise("us use user"), ShouldResemble, []Token{
			Token{Type: "identifier", Value: "us"},
			Token{Type: "use"},
			Token{Type: "identifier", Value: "user"},
		})
	})

	Convey("Import is its own token", t, func() {
		So(Tokenise("impor import imports"), ShouldResemble, []Token{
			Token{Type: "identifier", Value: "impor"},
			Token{Type: "import"},
			Token{Type: "identifier", Value: "imports"},
		})
	})

	Convey("When is its own token", t, func() {
		So(Tokenise("whe when whent"), ShouldResemble, []Token{
			Token{Type: "identifier", Value: "whe"},
			Token{Type: "when"},
			Token{Type: "identifier", Value: "whent"},
		})
	})

	Convey("Infix operators", t, func() {
		Convey("are found with spaces around them", func() {
			So(Tokenise("a eq b"), ShouldResemble, []Token{
				Token{Type: "identifier", Value: "a"},
				Token{Type: "infix_operator", Value: "eq"},
				Token{Type: "identifier", Value: "b"},
			})
		})

		Convey("without symbols are not found inside identifiers", func() {
			So(Tokenise("aeqb"), ShouldResemble, []Token{
				Token{Type: "identifier", Value: "aeqb"},
			})
		})

		Convey("with symbols are found inside identifiers", func() {
			So(Tokenise("a/b"), ShouldResemble, []Token{
				Token{Type: "identifier", Value: "a"},
				Token{Type: "infix_operator", Value: "/"},
				Token{Type: "identifier", Value: "b"},
			})
		})
	})

	Convey("Open block is found", t, func() {
		So(Tokenise("{"), ShouldResemble, []Token{
			Token{Type: "block_open"},
		})
	})

	Convey("Close block is found", t, func() {
		So(Tokenise("}"), ShouldResemble, []Token{
			Token{Type: "block_close"},
		})
	})

	Convey("Strings", t, func() {
		Convey("are found with double quotes", func() {
			So(Tokenise("\"Hello!\""), ShouldResemble, []Token{
				Token{Type: "string", Value: "Hello!"},
			})
		})

		Convey("can escape double quotes within them", nil)
	})

	Convey("Numbers", t, func() {

		Convey("can be integers", func() {
			So(Tokenise("1"), ShouldResemble, []Token{
				Token{Type: "number", Value: "1"},
			})
		})

		Convey("can be floats", func() {
			So(Tokenise("1.5"), ShouldResemble, []Token{
				Token{Type: "number", Value: "1.5"},
			})
		})

	})

	Convey("Booleans", t, func() {

		Convey("include true", func() {
			So(Tokenise("tru true truer"), ShouldResemble, []Token{
				Token{Type: "identifier", Value: "tru"},
				Token{Type: "boolean", Value: "true"},
				Token{Type: "identifier", Value: "truer"},
			})
		})

		Convey("include false", func() {
			So(Tokenise("fals false falser"), ShouldResemble, []Token{
				Token{Type: "identifier", Value: "fals"},
				Token{Type: "boolean", Value: "false"},
				Token{Type: "identifier", Value: "falser"},
			})
		})

	})

	Convey("Identifiers", t, func() {

		Convey("hashes", func() {

			Convey("can not start with a hash", func() {
				So(Tokenise("#id"), ShouldResemble, []Token{})
			})

			Convey("can not include a hash", func() {
				So(Tokenise("ab#cd"), ShouldResemble, []Token{
					Token{Type: "identifier", Value: "ab"},
				})
			})

		})

		Convey("brackets", func() {

			Convey("can not include an opening bracket", func() {
				So(Tokenise("ab(cd"), ShouldResemble, []Token{
					Token{Type: "identifier", Value: "ab"},
					Token{Type: "bracket_open"},
					Token{Type: "identifier", Value: "cd"},
				})
			})

			Convey("can not include a closing bracket", func() {
				So(Tokenise("ab)cd"), ShouldResemble, []Token{
					Token{Type: "identifier", Value: "ab"},
					Token{Type: "bracket_close"},
					Token{Type: "identifier", Value: "cd"},
				})
			})
		})

		Convey("can not include a comma", func() {
			So(Tokenise("ab,cd"), ShouldResemble, []Token{
				Token{Type: "identifier", Value: "ab"},
				Token{Type: "comma"},
				Token{Type: "identifier", Value: "cd"},
			})
		})

		Convey("can not include an end statement", func() {
			So(Tokenise("ab;cd"), ShouldResemble, []Token{
				Token{Type: "identifier", Value: "ab"},
				Token{Type: "end_statement"},
				Token{Type: "identifier", Value: "cd"},
			})
		})

		Convey("blocks", func() {
			Convey("can not include an opening block", func() {
				So(Tokenise("ab{cd"), ShouldResemble, []Token{
					Token{Type: "identifier", Value: "ab"},
					Token{Type: "block_open"},
					Token{Type: "identifier", Value: "cd"},
				})
			})

			Convey("can not include a closing block", func() {
				So(Tokenise("ab}cd"), ShouldResemble, []Token{
					Token{Type: "identifier", Value: "ab"},
					Token{Type: "block_close"},
					Token{Type: "identifier", Value: "cd"},
				})
			})
		})

		Convey("can not include a double quote", func() {
			So(Tokenise("ab\"cd"), ShouldResemble, []Token{
				Token{Type: "identifier", Value: "ab"},
				Token{Type: "string", Value: "cd"},
			})
		})

		Convey("can not include a space", func() {
			So(Tokenise("ab cd"), ShouldResemble, []Token{
				Token{Type: "identifier", Value: "ab"},
				Token{Type: "identifier", Value: "cd"},
			})
		})

		Convey("numbers", func() {
			Convey("can not start with a number", func() {
				So(Tokenise("1a"), ShouldResemble, []Token{
					Token{Type: "number", Value: "1"},
					Token{Type: "identifier", Value: "a"},
				})
			})

			Convey("can include a number", func() {
				So(Tokenise("a1"), ShouldResemble, []Token{
					Token{Type: "identifier", Value: "a1"},
				})
			})
		})

		Convey("symbols", func() {
			Convey("can start with symbols", func() {
				So(Tokenise("@"), ShouldResemble, []Token{
					Token{Type: "identifier", Value: "@"},
				})
			})

			Convey("can include symbols", func() {
				So(Tokenise("a@b"), ShouldResemble, []Token{
					Token{Type: "identifier", Value: "a@b"},
				})
			})
		})
	})
}
