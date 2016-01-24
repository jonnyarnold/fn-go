package parser

import (
	"github.com/jonnyarnold/fn-go/compiler/tokeniser"
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

	Convey("Expressions can be delimited by end_statement tokens", t, func() {
		exprs, err := Parse(tokensFor("foo; bar"))

		So(exprs, ShouldResemble, []Expression{
			IdentifierExpression{Name: "foo"},
			IdentifierExpression{Name: "bar"},
		})
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
			exprs, err := Parse(tokensFor("when{ true { foo } }"))

			So(exprs, ShouldResemble, []Expression{
				ConditionalExpression{
					Branches: []ConditionalBranchExpression{
						ConditionalBranchExpression{
							Condition: BooleanExpression{Value: true},
							Body: BlockExpression{
								Body: []Expression{
									IdentifierExpression{Name: "foo"},
								},
							},
						},
					},
				},
			})
			So(err, ShouldBeNil)
		})

	})

	Convey("Infix Operators", t, func() {

		Convey("become Function Call Expressions", func() {
			exprs, err := Parse(tokensFor("a + b"))

			So(exprs, ShouldResemble, []Expression{
				FunctionCallExpression{
					Identifier: IdentifierExpression{Name: "+"},
					Arguments: []Expression{
						IdentifierExpression{Name: "a"},
						IdentifierExpression{Name: "b"},
					},
				},
			})
			So(err, ShouldBeNil)
		})

		Convey("with the same precedence are executed left-to-right", func() {
			exprs, err := Parse(tokensFor("a * b * c"))

			So(exprs, ShouldResemble, []Expression{
				FunctionCallExpression{
					Identifier: IdentifierExpression{Name: "*"},
					Arguments: []Expression{
						IdentifierExpression{Name: "a"},
						FunctionCallExpression{
							Identifier: IdentifierExpression{Name: "*"},
							Arguments: []Expression{
								IdentifierExpression{Name: "b"},
								IdentifierExpression{Name: "c"},
							},
						},
					},
				},
			})
			So(err, ShouldBeNil)
		})

		Convey("with different precedence are executed in precedence order", func() {
			exprs, err := Parse(tokensFor("a + b * c"))

			So(exprs, ShouldResemble, []Expression{
				FunctionCallExpression{
					Identifier: IdentifierExpression{Name: "+"},
					Arguments: []Expression{
						IdentifierExpression{Name: "a"},
						FunctionCallExpression{
							Identifier: IdentifierExpression{Name: "*"},
							Arguments: []Expression{
								IdentifierExpression{Name: "b"},
								IdentifierExpression{Name: "c"},
							},
						},
					},
				},
			})
			So(err, ShouldBeNil)
		})

	})

	Convey("Blocks", t, func() {
		Convey("fail on non-primary statements", func() {
			_, err := Parse(tokensFor("{=}"))
			So(err, ShouldNotBeNil)
		})

		Convey("fail if the end of the tokens is reached before closing the block", func() {
			_, err := Parse(tokensFor("{foo"))
			So(err, ShouldNotBeNil)
		})
	})

	Convey("Brackets", t, func() {
		Convey("fail on non-primary statements", func() {
			_, err := Parse(tokensFor("(=)"))
			So(err, ShouldNotBeNil)
		})

		Convey("fail if the end of the tokens is reached before closing the bracket", func() {
			_, err := Parse(tokensFor("(foo"))
			So(err, ShouldNotBeNil)
		})

		Convey("provides precedence", func() {
			exprs, err := Parse(tokensFor("a * (b + c)"))

			So(exprs, ShouldResemble, []Expression{
				FunctionCallExpression{
					Identifier: IdentifierExpression{Name: "*"},
					Arguments: []Expression{
						IdentifierExpression{Name: "a"},
						FunctionCallExpression{
							Identifier: IdentifierExpression{Name: "+"},
							Arguments: []Expression{
								IdentifierExpression{Name: "b"},
								IdentifierExpression{Name: "c"},
							},
						},
					},
				},
			})
			So(err, ShouldBeNil)
		})

	})

	Convey("Function Calls", t, func() {
		Convey("fail on non-value parameters", func() {
			_, err := Parse(tokensFor("foo(=)"))
			So(err, ShouldNotBeNil)
		})

		Convey("fail on an invalid parameter list", func() {
			_, err := Parse(tokensFor("foo(a b)"))
			So(err, ShouldNotBeNil)
		})

		Convey("fail if the parameter list is never closed", func() {
			_, err := Parse(tokensFor("foo(a"))
			So(err, ShouldNotBeNil)
		})

		Convey("permit valid parameter lists", func() {
			_, err := Parse(tokensFor("foo(a, b)"))
			So(err, ShouldBeNil)
		})
	})

	Convey("Function Definitions", t, func() {
		Convey("fail on an invalid argument list", func() {
			_, err := Parse(tokensFor("foo(a b) { true }"))
			So(err, ShouldNotBeNil)
		})

		Convey("fail if the argument list is never closed", func() {
			_, err := Parse(tokensFor("foo(a { true }"))
			So(err, ShouldNotBeNil)
		})

		Convey("fail if arguments are not identifiers", func() {
			_, err := Parse(tokensFor("foo(a, =) { true }"))
			So(err, ShouldNotBeNil)
		})

		Convey("permit valid argument lists", func() {
			_, err := Parse(tokensFor("foo(a, b) { true }"))
			So(err, ShouldBeNil)
		})
	})

	Convey("Conditionals", t, func() {

		Convey("fail if a block is not opened after the 'when'", func() {
			_, err := Parse(tokensFor("when foo"))
			So(err, ShouldNotBeNil)
		})

		Convey("fail if a non-value is passed as a condition", func() {
			_, err := Parse(tokensFor("when { = { true } }"))
			So(err, ShouldNotBeNil)
		})

		Convey("fail if a block is not passed after a condition", func() {
			_, err := Parse(tokensFor("when { foo }"))
			So(err, ShouldNotBeNil)
		})

		Convey("can have multiple branches", func() {
			exprs, err := Parse(tokensFor("when { false { foo } true { bar } }"))
			So(exprs, ShouldResemble, []Expression{
				ConditionalExpression{
					Branches: []ConditionalBranchExpression{
						ConditionalBranchExpression{
							Condition: BooleanExpression{Value: false},
							Body: BlockExpression{
								Body: []Expression{
									IdentifierExpression{Name: "foo"},
								},
							},
						},
						ConditionalBranchExpression{
							Condition: BooleanExpression{Value: true},
							Body: BlockExpression{
								Body: []Expression{
									IdentifierExpression{Name: "bar"},
								},
							},
						},
					},
				},
			})
			So(err, ShouldBeNil)
		})

	})
}
