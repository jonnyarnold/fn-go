package vm

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestExecute(t *testing.T) {
	Convey("Machine", t, func() {

		Convey("is not running by default", nil)

		Convey("keeps state once run", nil)

		Convey("raises an error if no instructions are passed", nil)

		Convey("raises no error if only TERMINATE is passed", nil)

		Convey("DECLARE_CONSTANT", nil)

		Convey("SET_REGISTER_WITH_CONSTANT", nil)

		Convey("ADD", nil)

		Convey("SUBTRACT", nil)

		Convey("MULTIPLY", nil)

		Convey("DIVIDE", nil)

		Convey("AND", nil)

		Convey("OR", nil)

		Convey("DUMP", nil)

	})
}
