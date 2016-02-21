package vm

import (
	"github.com/jonnyarnold/fn-go/bytecode"
)

func RunBytecode() string {
	// x = (y) {
	//   when {
	//     y { not(y) }
	//     true { y }
	//   }
	// }
	// x(false)
	// x(true)
	code := bytecode.New().
		FunctionFollows(
		bytecode.New().
			JumpIfFalse(bytecode.New().
			Negate()).
			Return()).
		Dump().
		DeclareBool(false).
		Call(0).
		Dump().
		DeclareBool(true).
		Call(0).
		Dump().
		Return()

	machine := NewMachine(code)
	machine.ProcessAll()
	return machine.Result()
}
