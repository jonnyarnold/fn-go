package parser

import (
	"github.com/jonnyarnold/fn-go/tokeniser"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func tokensFor(code string) []tokeniser.Token {
	return tokeniser.Tokenise(code)
}

func TestParser(t *testing.T) {
	Convey("An empty token list gives an empty expression list", t, func() {
		exprs, err := Parse(tokensFor(""))

		So(exprs, ShouldResemble, []Expression{})
		So(err, ShouldBeNil)
	})

	Convey("Parsing value", t, func() {

		Convey("identifiers become Identifier Expressions", func() {
			exprs, err := Parse(tokensFor("foo"))

			So(exprs, ShouldResemble, []Expression{
				IdentifierExpression{Name: "foo"},
			})
			So(err, ShouldBeNil)
		})

		Convey("function calls become Function Call Expressions", func() {
			exprs, err := Parse(tokensFor("foo(bar)"))

			So(exprs, ShouldResemble, []Expression{
				FunctionCallExpression{
					Identifier: IdentifierExpression{Name: "foo"},
					Arguments: []Expression{
						IdentifierExpression{Name: "bar"},
					},
				},
			})
			So(err, ShouldBeNil)
		})

		Convey("numbers become Number Expressions", func() {
			exprs, err := Parse(tokensFor("1.234"))

			So(exprs, ShouldResemble, []Expression{
				NumberExpression{Value: "1.234"},
			})
			So(err, ShouldBeNil)
		})

		Convey("strings become String Expressions", func() {
			exprs, err := Parse(tokensFor("\"foo\""))

			So(exprs, ShouldResemble, []Expression{
				StringExpression{Value: "foo"},
			})
			So(err, ShouldBeNil)
		})

		Convey("booleans become Boolean Expressions", func() {
			exprs, err := Parse(tokensFor("true"))

			So(exprs, ShouldResemble, []Expression{
				BooleanExpression{Value: true},
			})
			So(err, ShouldBeNil)
		})

		Convey("brackets enclose other primary expressions", func() {
			exprs, err := Parse(tokensFor("(foo)"))

			So(exprs, ShouldResemble, []Expression{
				IdentifierExpression{Name: "foo"},
			})
			So(err, ShouldBeNil)
		})

		Convey("function definitions become Function Definition Expressions", func() {
			exprs, err := Parse(tokensFor("(a) { a }"))

			So(exprs, ShouldResemble, []Expression{
				FunctionPrototypeExpression{
					Arguments: []IdentifierExpression{
						IdentifierExpression{Name: "a"},
					},
					Body: BlockExpression{
						Body: []Expression{
							IdentifierExpression{Name: "a"},
						},
					},
				},
			})
			So(err, ShouldBeNil)
		})

		Convey("fail on an unexpected token", func() {
			_, err := Parse(tokensFor("="))
			So(err, ShouldNotBeNil)
		})

		Convey("blocks become Block Expressions", func() {
			exprs, err := Parse(tokensFor("{ a }"))

			So(exprs, ShouldResemble, []Expression{
				BlockExpression{
					Body: []Expression{
						IdentifierExpression{Name: "a"},
					},
				},
			})
			So(err, ShouldBeNil)
		})

		Convey("when statements become Conditional Expressions", func() {
			exprs, err := Parse(tokensFor("when{}"))

			So(exprs, ShouldResemble, []Expression{
				ConditionalExpression{},
			})
			So(err, ShouldBeNil)
		})

	})

	Convey("Infix Operators", t, func() {

		Convey("become Function Call Expressions", nil)

		Convey("with the same precedence are executed left-to-right", nil)

		Convey("with different precedence are executed in precedence order", nil)

	})

	Convey("Blocks", t, func() {

		Convey("fail on non-primary statements", nil)

		Convey("fail if the end of the tokens is reached before closing the block", nil)

	})

	Convey("Brackets", t, func() {

		Convey("enclose primary statements", nil)

		Convey("fail on non-primary statements", nil)

		Convey("fail if the end of the tokens is reached before closing the brackets", nil)

		Convey("provides precedence", nil)

	})

	Convey("Function Calls", t, func() {

		Convey("fail on non-value parameters", nil)

		Convey("fail on an invalid parameter list", nil)

		Convey("fail if the parameter list is never closed", nil)

	})

	Convey("Function Definitions", t, func() {

		Convey("fail on non-identifier arguments", nil)

		Convey("fail on an invalid argument list", nil)

		Convey("fail if the argument list is never closed", nil)

	})

	Convey("Conditionals", t, func() {

		Convey("fail if a block is not opened after the 'when'", nil)

		Convey("fail if a non-value is passed as a condition", nil)

		Convey("fail if a block is not passed after a condition", nil)

		Convey("can have multiple branches", nil)

	})
}
