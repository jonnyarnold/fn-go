package runtime

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestExecute(t *testing.T) {
	Convey("An empty expression list", t, func() {

		Convey("returns no value", nil)

		Convey("returns no error", nil)

	})

	Convey("Number expressions return numeric values", t, nil)

	Convey("String expressions return string values", t, nil)

	Convey("Boolean expressions return boolean values", t, nil)

	Convey("Identifier expressions", t, func() {

		Convey("return an error if not defined", nil)

		Convey("return the value if defined", nil)

	})

	Convey("Function prototype expressions", t, func() {

		Convey("create a new child scope", nil)

		Convey("return a Go func", nil)

	})

	Convey("Block expressions", t, func() {

		Convey("can redefine parent definitions", nil)

		Convey("return a scope with definitions applied", nil)

	})

	Convey("Function call expressions", t, func() {

		Convey("return an error if not defined", nil)

		Convey("executes arguments", nil)

		Convey("returns the function value", nil)

		Convey("=", func() {

			Convey("sets the given ID to the given value", nil)

			Convey("returns the new scope", nil)

			Convey("does not return a value", nil)

			Convey("returns an error if ID already defined", nil)

		})

		Convey(".", func() {

			Convey("executes the child in the parent scope", nil)

			Convey("returns an error if the parent scope is not defined", nil)

			Convey("returns the grandparent scope", nil)

			Convey("returns the value from the child expression", nil)

		})

		Convey("for user-defined functions", func() {

			Convey("return the value of the prototype", nil)

			Convey("return an error on mismatched argument lengths", nil)

			Convey("execute in the function prototype scope", nil)

		})

	})

	Convey("Conditional expressions", t, func() {

		Convey("return the block's value if the condition is true", nil)

		Convey("return an error if no conditions are met", nil)

		Convey("stops executing conditions on the first true condition", nil)

	})

	Convey("Built-in functions", t, func() {

		Convey("are available from the expression list", nil)

		Convey("are available within user-defined functions", nil)

		Convey("are available within user-defined blocks", nil)

		Convey("are available within conditions", nil)

		Convey("+", func() {

			Convey("sums two numbers", nil)

		})

		Convey("-", func() {

			Convey("takes the difference of two numbers", nil)

		})

		Convey("*", func() {

			Convey("multiplies two numbers", nil)

		})

		Convey("/", func() {

			Convey("divides two numbers", nil)

		})

		Convey("and", func() {

			Convey("returns true with two true boolean values", nil)

			Convey("returns false if one value is false", nil)

			Convey("does not execute the second argument if the first is false", nil)

		})

		Convey("or", func() {

			Convey("returns true if one value is true", nil)

			Convey("returns false with two false values", nil)

			Convey("does not execute the second argument if the first is true", nil)

		})

		Convey("not", func() {

			Convey("returns false when the condition is true", nil)

			Convey("returns true when the condition is false", nil)

		})

		Convey("eq", func() {

			Convey("returns true if two numbers have the same value", nil)

			Convey("returns false if two numbers are different", nil)

			Convey("returns true if two strings have the same value", nil)

			Convey("returns false if two strings are different", nil)

			Convey("returns true if two booleans have the same value", nil)

			Convey("returns false if two booleans are different", nil)

		})

		Convey("print", func() {

			Convey("outputs to the console", nil)

			Convey("returns no value", nil)

		})

		Convey("List", func() {

			Convey("returns a List with the given arguments", nil)

		})

	})
}
