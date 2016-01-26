package vm

import (
	"fmt"
)

// The number of bytes in an instruction.
// Maximum is:
// DECLARE_CONSTANT(1) INT(1) INT8(1)
const INSTRUCTION_BYTES = 3

// The number of registers in the Machine.
const NUM_REGISTERS = 2

// The Machine is the Virtual Machine
// used to make calculations.
type Machine struct {
	// Code is used to store the bytecode.
	Code []byte

	// Register is an array of registers, mutable constantValue objects.
	Register [NUM_REGISTERS]constantValue

	// Constant is an array of constants.
	Constant []constantValue

	// InstructionCounter tracks the current or last
	// processed instruction.
	//
	// Before execution, this is 0.
	// The first instruction is instruction 1.
	InstructionCounter uint

	// Running is set to true when started, and is used to
	// break out of the instruction loop when complete.
	Running bool
}

// Initializes a new Machine.
func NewMachine(bytecode []byte) Machine {
	return Machine{
		Code:    bytecode,
		Running: false,
	}
}

// Processes all instructions in Code
// until the machines stops Running.
func (m *Machine) ProcessAll() {
	m.Running = true

	for m.Running {
		m.ProcessNext()
	}
}

// Process the next instruction in the Code.
func (m *Machine) ProcessNext() {
	instruction := m.GetNextInstruction()
	m.ExecuteInstruction(instruction)
}

// Get the next instruction in the Code.
//
// This also advances the InstructionCounter.
func (m *Machine) GetNextInstruction() []byte {
	nextStartIdx := m.InstructionCounter * INSTRUCTION_BYTES
	nextEndIdx := nextStartIdx + INSTRUCTION_BYTES

	fmt.Printf("Getting [%d:%d]: ", nextStartIdx, nextEndIdx)
	instruction := m.Code[nextStartIdx:nextEndIdx]
	fmt.Printf("%x\n", instruction)

	// This needs to be after we find the instruction,
	// otherwise we'll never read instruction 0.
	m.InstructionCounter += 1

	return instruction
}

func (m *Machine) ExecuteInstruction(instruction []byte) {
	// The first byte is the instruction type.
	switch instruction[0] {

	// TERMINATE
	// Stops program execution.
	case TERMINATE:
		m.Running = false
		break

	// DECLARE_CONSTANT constant
	// Adds the constant to the Constant List.
	case DECLARE_CONSTANT:
		m.Constant = append(m.Constant, ConstantValue(instruction[1:3]))
		break

	// SET_REGISTER_WITH_CONSTANT register_idx constant_idx
	// Pushes the constant into the register.
	case SET_REGISTER_WITH_CONSTANT:
		m.Register[instruction[1]] = m.Constant[instruction[2]]
		break
	}
}

func (m *Machine) Result() string {
	return fmt.Sprintf("%v", m)
}

func RunBytecode() string {
	bytecode := []byte{
		DECLARE_CONSTANT, TYPE_INT, 0x01,
		DECLARE_CONSTANT, TYPE_BOOL, 0x01,
		TERMINATE, 0x00, 0x00,
	}

	machine := NewMachine(bytecode)
	machine.ProcessAll()
	return machine.Result()
}
