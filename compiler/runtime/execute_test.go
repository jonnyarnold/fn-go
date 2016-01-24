package runtime

import (
	. "github.com/jonnyarnold/fn-go/compiler/parser"
	"github.com/jonnyarnold/fn-go/compiler/tokeniser"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

// Returns the Expressions for the given code.
func exprsFor(code string) []Expression {
	tokens := tokeniser.Tokenise(code)
	exprs, err := Parse(tokens)

	if err != nil {
		panic(err)
	}

	return exprs
}

func eval(code string) EvalResult {
	return Execute(exprsFor(code))
}

func TestExecute(t *testing.T) {
	Convey("An empty expression list", t, func() {

		Convey("returns no value", func() {
			result := eval("")
			So(result.Value, ShouldBeNil)
		})

		Convey("returns no error", func() {
			result := eval("")
			So(result.Error, ShouldBeNil)
		})

	})

	Convey("Number expressions", t, func() {

		Convey("return numeric values", func() {
			result := eval("2.5")

			So(result.Value, ShouldResemble, number{value: "2.5"})
			So(result.Error, ShouldBeNil)
		})

		Convey("cannot be extended", func() {
			result := eval("(1).foo = 1")

			So(result.Error, ShouldNotBeNil)
		})

	})

	Convey("String expressions", t, func() {

		Convey("return string values", func() {
			result := eval("\"Hi\"")

			So(result.Value, ShouldResemble, fnString{value: "Hi"})
			So(result.Error, ShouldBeNil)
		})

		Convey("cannot be extended", func() {
			result := eval("(\"Hi\").foo = 1")

			So(result.Error, ShouldNotBeNil)
		})

	})

	Convey("Boolean expressions", t, func() {

		Convey("return boolean values", func() {
			result := eval("false")

			So(result.Value, ShouldResemble, fnBool{value: false})
			So(result.Error, ShouldBeNil)
		})

		Convey("cannot be extended", func() {
			result := eval("(true).foo = 1")

			So(result.Error, ShouldNotBeNil)
		})

	})

	Convey("Identifier expressions", t, func() {

		Convey("return an error if not defined", func() {
			result := eval("notDefined")

			So(result.Error, ShouldNotBeNil)
		})

		Convey("return the value if defined", func() {
			result := eval("print")

			So(result.Value, ShouldNotBeNil)
			So(result.Error, ShouldBeNil)
		})

	})

	Convey("Function prototype expressions", t, func() {

		Convey("create a new child scope", nil)

		Convey("return a functionScope", func() {
			result := eval("() { }")

			So(result.Value, ShouldHaveSameTypeAs, functionScope{})
			So(result.Error, ShouldBeNil)
		})

		Convey("can return blocks", func() {
			result := eval("x = () { y = \"foo\" }; z = x(); z.y")

			So(result.Error, ShouldBeNil)
			So(result.Value, ShouldResemble, fnString{value: "foo"})
		})

	})

	Convey("Block expressions", t, func() {

		Convey("can redefine parent definitions", func() {
			result := eval("x = { print = \"Hello\" }; x.print")

			So(result.Value, ShouldResemble, fnString{value: "Hello"})
			So(result.Error, ShouldBeNil)
		})

		Convey("return a scope with definitions applied", func() {
			result := eval("{ foo = \"bar\" }")

			So(result.Value, ShouldHaveSameTypeAs, Scope{})
			So(result.Value.(Scope).definitions, ShouldResemble, defMap{
				"foo": fnString{value: "bar"},
			})
		})

		Convey("can be called if a call attribute is defined", func() {
			result := eval("x = { call = (a) { a } }; x(1)")

			So(result.Error, ShouldBeNil)
			So(result.Value, ShouldResemble, number{value: "1"})
		})

		Convey("return an error if a call attribute is not defined and is called", func() {
			result := eval("x = { }; x(1)")
			So(result.Error, ShouldNotBeNil)
		})

	})

	Convey("Function call expressions", t, func() {

		Convey("return an error if not defined", func() {
			result := eval("notDefined()")

			So(result.Error, ShouldNotBeNil)
		})

		Convey("returns the function value", func() {
			result := eval("returnX = (x) { x }; returnX(1)")

			So(result.Value, ShouldResemble, number{value: "1"})
			So(result.Error, ShouldBeNil)
		})

		Convey("executes arguments", func() {
			result := eval("returnX = (x) { x }; returnX(2 + 2)")

			So(result.Error, ShouldBeNil)
			So(result.Value, ShouldResemble, number{value: "4"})
		})

		Convey("=", func() {

			Convey("sets the given ID to the given value", func() {
				result := eval("x = 2; x")

				So(result.Error, ShouldBeNil)
				So(result.Value, ShouldResemble, number{value: "2"})
			})

			Convey("returns the defining block", func() {
				result := eval("foo = \"bar\"")

				So(result.Error, ShouldBeNil)
				So(result.Value.(Scope).definitions, ShouldResemble, defMap{
					"foo": fnString{value: "bar"},
				})
			})

			Convey("returns an error if ID already defined", func() {
				result := eval("x = 1; x = 2")

				So(result.Error, ShouldNotBeNil)
			})

		})

		Convey(".", func() {

			Convey("executes the child in the parent scope", func() {
				result := eval("x = { print = (x) { x } }; x.print(1)")

				So(result.Error, ShouldBeNil)
				So(result.Value, ShouldResemble, number{value: "1"})
			})

			Convey("returns an error if the parent scope is not defined", func() {
				result := eval("x.print(1)")

				So(result.Error, ShouldNotBeNil)
			})

		})

		Convey("import!", func() {

			Convey("includes the file content directly into the current scope", func() {
				result := eval("import!(\"test_import.fn\"); x")

				So(result.Error, ShouldBeNil)
				So(result.Value, ShouldResemble, number{value: "1"})
			})

			Convey("returns an error if the file does not exist", func() {
				result := eval("import!(\"DOES-NOT-EXIST\")")

				So(result.Error, ShouldNotBeNil)
			})

			Convey("returns an error if variables are already defined", func() {
				result := eval("x = 1; import!(\"test_import.fn\")")

				So(result.Error, ShouldNotBeNil)
			})

		})

		Convey("import", func() {

			Convey("includes the file content in the variable", func() {
				result := eval("a = import(\"test_import.fn\"); a.x")

				So(result.Error, ShouldBeNil)
				So(result.Value, ShouldResemble, number{value: "1"})
			})

			Convey("returns an error if the file does not exist", func() {
				result := eval("a = import(\"DOES-NOT-EXIST\")")

				So(result.Error, ShouldNotBeNil)
			})

		})

		Convey("for user-defined functions", func() {

			Convey("return an error on mismatched argument lengths", func() {
				result := eval("x = (a) { a }; x(1, 2)")

				So(result.Error, ShouldNotBeNil)
			})

			Convey("return an error on argument errors", func() {
				result := eval("x = (a) { a }; x(foo)")

				So(result.Error, ShouldNotBeNil)
			})

			Convey("execute in the function prototype scope", nil) // func() {
			// result := eval("x = (a) { print = a; print }; x(1)")

			// So(result.Error, ShouldBeNil)
			// So(result.Value, ShouldResemble, number{value: "1"})
			// })

		})

	})

	Convey("Conditional expressions", t, func() {

		Convey("return the block's value if the condition is true", func() {
			result := eval("when { true { 1 } }")

			So(result.Error, ShouldBeNil)
			So(result.Value, ShouldResemble, number{value: "1"})
		})

		Convey("return an error if no conditions are met", func() {
			result := eval("when { false { 1 } }")

			So(result.Error, ShouldNotBeNil)
		})

		Convey("return an error if conditions contain errors", func() {
			result := eval("when { foo { 1 } }")

			So(result.Error, ShouldNotBeNil)
		})

		Convey("stops executing conditions on the first true condition", func() {
			result := eval("when { true { 1 } true { 2 } }")

			So(result.Error, ShouldBeNil)
			So(result.Value, ShouldResemble, number{value: "1"})
		})

	})

	Convey("Built-in functions", t, func() {

		Convey("are available from the expression list", func() {
			result := eval("print")

			So(result.Error, ShouldBeNil)
			So(result.Value, ShouldHaveSameTypeAs, functionScope{})
		})

		Convey("are available within user-defined functions", func() {
			result := eval("x = (a) { print(a) }; x(1)")

			So(result.Error, ShouldBeNil)
		})

		Convey("cannot be extended", func() {
			result := eval("print.foo = 1")

			So(result.Error, ShouldNotBeNil)
		})

		Convey("are available within user-defined blocks", func() {
			result := eval("x = { foo = print }; x.foo(1)")

			So(result.Error, ShouldBeNil)
		})

		Convey("are available within conditions", func() {
			result := eval("when { true eq true { 1 } }")

			So(result.Error, ShouldBeNil)
		})

		Convey("+", func() {

			Convey("sums two integers", func() {
				result := eval("2 + 2")
				So(result.Error, ShouldBeNil)
				So(result.Value, ShouldResemble, number{value: "4"})
			})

			Convey("sums two floats", func() {
				result := eval("2.5 + 2.5")
				So(result.Error, ShouldBeNil)
				So(result.Value, ShouldResemble, number{value: "5"})
			})

			Convey("sums an integer and a float", func() {
				result := eval("2.5 + 2")
				So(result.Error, ShouldBeNil)
				So(result.Value, ShouldResemble, number{value: "4.5"})
			})

		})

		Convey("-", func() {

			Convey("takes the difference of two integers", func() {
				result := eval("4 - 3")
				So(result.Error, ShouldBeNil)
				So(result.Value, ShouldResemble, number{value: "1"})
			})

			Convey("takes the difference of two floats", func() {
				result := eval("4.5 - 3.5")
				So(result.Error, ShouldBeNil)
				So(result.Value, ShouldResemble, number{value: "1"})
			})

			Convey("takes the difference of an integer and a float", func() {
				result := eval("4.5 - 3")
				So(result.Error, ShouldBeNil)
				So(result.Value, ShouldResemble, number{value: "1.5"})
			})

		})

		Convey("*", func() {

			Convey("multiplies two integers", func() {
				result := eval("2 * 2")
				So(result.Error, ShouldBeNil)
				So(result.Value, ShouldResemble, number{value: "4"})
			})

			Convey("multiplies two floats", func() {
				result := eval("1.5 * 1.5")
				So(result.Error, ShouldBeNil)
				So(result.Value, ShouldResemble, number{value: "2.25"})
			})

			Convey("multiplies an integer and a float", func() {
				result := eval("2 * 1.5")
				So(result.Error, ShouldBeNil)
				So(result.Value, ShouldResemble, number{value: "3"})
			})

		})

		Convey("/", func() {

			Convey("divides two integers", func() {
				result := eval("4 / 2")
				So(result.Error, ShouldBeNil)
				So(result.Value, ShouldResemble, number{value: "2"})
			})

			Convey("divides two floats", func() {
				result := eval("2.5 / 2.5")
				So(result.Error, ShouldBeNil)
				So(result.Value, ShouldResemble, number{value: "1"})
			})

			Convey("divides an integer and a float", func() {
				result := eval("2.5 / 5")
				So(result.Error, ShouldBeNil)
				So(result.Value, ShouldResemble, number{value: "0.5"})
			})

		})

		Convey("and", func() {

			Convey("returns true with two true boolean values", func() {
				result := eval("true and true")
				So(result.Error, ShouldBeNil)
				So(result.Value, ShouldResemble, fnBool{value: true})
			})

			Convey("returns false if one value is false", func() {
				result := eval("false and true")
				So(result.Error, ShouldBeNil)
				So(result.Value, ShouldResemble, fnBool{value: false})
			})

		})

		Convey("or", func() {

			Convey("returns true if one value is true", func() {
				result := eval("true or false")
				So(result.Error, ShouldBeNil)
				So(result.Value, ShouldResemble, fnBool{value: true})
			})

			Convey("returns false with two false values", func() {
				result := eval("false or false")
				So(result.Error, ShouldBeNil)
				So(result.Value, ShouldResemble, fnBool{value: false})
			})

		})

		Convey("not", func() {

			Convey("returns false when the condition is true", func() {
				result := eval("not(true)")
				So(result.Error, ShouldBeNil)
				So(result.Value, ShouldResemble, fnBool{value: false})
			})

			Convey("returns true when the condition is false", func() {
				result := eval("not(false)")
				So(result.Error, ShouldBeNil)
				So(result.Value, ShouldResemble, fnBool{value: true})
			})

		})

		Convey("eq", func() {

			Convey("returns true if two numbers have the same value", func() {
				result := eval("2 eq 2")
				So(result.Error, ShouldBeNil)
				So(result.Value, ShouldResemble, fnBool{value: true})
			})

			Convey("returns false if two numbers are different", func() {
				result := eval("2 eq 1")
				So(result.Error, ShouldBeNil)
				So(result.Value, ShouldResemble, fnBool{value: false})
			})

			Convey("returns true if two strings have the same value", func() {
				result := eval("\"foo\" eq \"foo\"")
				So(result.Error, ShouldBeNil)
				So(result.Value, ShouldResemble, fnBool{value: true})
			})

			Convey("returns false if two strings are different", func() {
				result := eval("\"foo\" eq \"bar\"")
				So(result.Error, ShouldBeNil)
				So(result.Value, ShouldResemble, fnBool{value: false})
			})

			Convey("returns true if two booleans have the same value", func() {
				result := eval("true eq true")
				So(result.Error, ShouldBeNil)
				So(result.Value, ShouldResemble, fnBool{value: true})
			})

			Convey("returns false if two booleans are different", func() {
				result := eval("true eq false")
				So(result.Error, ShouldBeNil)
				So(result.Value, ShouldResemble, fnBool{value: false})
			})

		})

		Convey("print", func() {

			Convey("outputs to the console", func() {
				result := eval("print(1)")
				So(result.Error, ShouldBeNil)
			})

			Convey("returns no value", func() {
				result := eval("print(1)")
				So(result.Error, ShouldBeNil)
				So(result.Value, ShouldBeNil)
			})

		})

		Convey("List", func() {

			Convey("returns a List with the given arguments", func() {
				result := eval("List(1, \"two\", true)")
				So(result.Error, ShouldBeNil)
				So(result.Value, ShouldResemble, list{
					Items: []fnScope{
						number{value: "1"},
						fnString{value: "two"},
						fnBool{value: true},
					},
				})
			})

		})

	})
}
