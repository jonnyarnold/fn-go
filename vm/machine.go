package vm

import (
	"fmt"
)

// The number of bytes in an instruction.
const INSTRUCTION_BYTES = 4

// The number of addressable registers.
const NUM_REGISTERS = 4

// The number of addressable constants.
const NUM_CONSTANTS = 4

// The Machine is the Virtual Machine
// used to make calculations.
type Machine struct {
	// Code is used to store the bytecode.
	Code []byte

	// Register is an array of registers, mutable constantValue objects.
	Register [NUM_REGISTERS]vmConstant

	// Constant is an array of constants.
	Constant []vmConstant

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
		m.Constant = append(m.Constant, VMConstant(instruction[1:3]))
		break

	// SET_REGISTER_WITH_CONSTANT register_idx constant_idx
	// Pushes the constant into the register.
	case SET_REGISTER_WITH_CONSTANT:
		m.Register[instruction[1]] = m.Constant[instruction[2]]
		break

	// ADD result_register_id operand_register_1_id operand_register_2_id
	// Adds the values of the two operand registers together and places them in the result register.
	case ADD:
		m.Register[instruction[1]] = AddNumbers(m.Register[instruction[2]].(vmNumber), m.Register[instruction[3]].(vmNumber))
		break

	// SUBTRACT result_register_id operand_register_1_id operand_register_2_id
	// Subtracts the values of the two operand registers together and places them in the result register.
	case SUBTRACT:
		m.Register[instruction[1]] = SubtractNumbers(m.Register[instruction[2]].(vmNumber), m.Register[instruction[3]].(vmNumber))
		break

	// MULTIPLY result_register_id operand_register_1_id operand_register_2_id
	// Multiplies the values of the two operand registers together and places them in the result register.
	case MULTIPLY:
		m.Register[instruction[1]] = MultiplyNumbers(m.Register[instruction[2]].(vmNumber), m.Register[instruction[3]].(vmNumber))
		break

	// DIVIDE result_register_id operand_register_1_id operand_register_2_id
	// Divides the values of the two operand registers together and places them in the result register.
	case DIVIDE:
		m.Register[instruction[1]] = DivideNumbers(m.Register[instruction[2]].(vmNumber), m.Register[instruction[3]].(vmNumber))
		break

	// AND result_register_id operand_register_1_id operand_register_2_id
	// Takes the logical AND of the two operand registers.
	case AND:
		m.Register[instruction[1]] = m.Register[instruction[2]].(vmBool).Value && m.Register[instruction[3]].(vmBool).Value
		break

	// OR result_register_id operand_register_1_id operand_register_2_id
	// Takes the logical OR of the two operand registers.
	case OR:
		m.Register[instruction[1]] = m.Register[instruction[2]].(vmBool).Value || m.Register[instruction[3]].(vmBool).Value
		break

	// DUMP
	// Dumps the current machine state
	case DUMP:
		fmt.Println(m)
	}

}

func (m *Machine) Result() string {
	return fmt.Sprintf("%v", m)
}

func (m *Machine) String() string {
	return fmt.Sprintf("VM STATE:\nRegisters: %v\nConstants: %v", m.Register, m.Constant)
}
