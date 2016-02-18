package vm

func RunBytecode() string {
	// x = (y) { when { y { not(y) } true { y } } }
	// x(false)
	// x(true)
	bytecode := []byte{
		DECLARE, TYPE_BOOL, 0x00,
		DECLARE, TYPE_BOOL, 0x01,
		FUNCTION_FOLLOWS, 0x04, 0x00,
		NEGATE, 0x00, 0x00,
		JUMP_IF_FALSE, 0x02, 0x00,
		NEGATE, 0x00, 0x00,
		RETURN, 0x00, 0x00,
		COPY, 0x00, 0x00,
		CALL, 0x02, 0x00,
		COPY, 0x01, 0x00,
		CALL, 0x02, 0x00,
		RETURN, 0x00, 0x00,
	}

	machine := NewMachine(bytecode)
	machine.ProcessAll()
	return machine.Result()
}
