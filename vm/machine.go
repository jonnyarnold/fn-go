package vm

import (
	"fmt"
)

// The number of bytes in an instruction.
const INSTRUCTION_BYTES = 3

// The Machine is the Virtual Machine
// used to make calculations.
type Machine struct {
	// Code is used to store the bytecode.
	Code []byte

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

	// The set of return points.
	// When a RETURN instruction is met,
	// the top of this stack replaces the instruction counter.
	Returns []uint
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

	instruction := m.Code[nextStartIdx:nextEndIdx]

	return instruction
}

func (m *Machine) ExecuteInstruction(instruction []byte) {
	// The first byte is the instruction type.
	switch instruction[0] {

	// DECLARE constant
	// Adds the constant to the Constant List.
	case DECLARE:
		fmt.Printf("[%v] DECLARE %v\n", m.InstructionCounter, instruction[1:3])
		m.Constant = append(m.Constant, VMConstant(instruction[1:3]))
		m.InstructionCounter += 1
		break

	// ADD
	// Adds the top 2 constants together and pushes them onto the Constant List.
	case ADD:
		m.Constant = append(
			m.Constant,
			AddNumbers(m.Constant[len(m.Constant)-1].(vmNumber), m.Constant[len(m.Constant)-2].(vmNumber)))
		m.InstructionCounter += 1
		break

	// SUBTRACT
	// Subtracts the top 2 constants together and pushes them onto the Constant List.
	case SUBTRACT:
		m.Constant = append(
			m.Constant,
			SubtractNumbers(m.Constant[len(m.Constant)-1].(vmNumber), m.Constant[len(m.Constant)-2].(vmNumber)))
		m.InstructionCounter += 1
		break

	// MULTIPLY
	// Multiplies the top 2 constants together and pushes them onto the Constant List.
	case MULTIPLY:
		m.Constant = append(
			m.Constant,
			MultiplyNumbers(m.Constant[len(m.Constant)-1].(vmNumber), m.Constant[len(m.Constant)-2].(vmNumber)))
		m.InstructionCounter += 1
		break

	// DIVIDE
	// Divides the top 2 constants together and pushes them onto the Constant List.
	case DIVIDE:
		m.Constant = append(
			m.Constant,
			DivideNumbers(m.Constant[len(m.Constant)-1].(vmNumber), m.Constant[len(m.Constant)-2].(vmNumber)))
		m.InstructionCounter += 1
		break

	// AND
	// Takes the logical AND of the top 2 constants.
	case AND:
		and := m.Constant[len(m.Constant)-1].(vmBool).Value && m.Constant[len(m.Constant)-2].(vmBool).Value
		m.Constant = append(
			m.Constant,
			vmBool{Value: and})
		m.InstructionCounter += 1
		break

	// OR
	// Takes the logical OR of the top 2 constants.
	case OR:
		or := m.Constant[len(m.Constant)-1].(vmBool).Value || m.Constant[len(m.Constant)-2].(vmBool).Value
		m.Constant = append(
			m.Constant,
			vmBool{Value: or})
		m.InstructionCounter += 1
		break

	// NEGATE
	// Takes the logical NOT or unary - of the top constant.
	case NEGATE:
		fmt.Printf("[%v] NEGATE\n", m.InstructionCounter)
		m.Constant = append(
			m.Constant,
			m.Constant[len(m.Constant)-1].Negate())
		m.InstructionCounter += 1
		break

	// DUMP
	// Dumps the current machine state
	case DUMP:
		fmt.Printf("[%v] DUMP\n", m.InstructionCounter)
		fmt.Println(m)
		m.InstructionCounter += 1
		break

	// RETURN
	// If there is nothing on the Return list, terminate.
	// If there is a Return, pop the top return from the stack
	// and go to it.
	case RETURN:
		fmt.Printf("[%v] RETURN\n", m.InstructionCounter)
		if len(m.Returns) == 0 {
			m.Running = false
		} else {
			m.InstructionCounter = m.Returns[len(m.Returns)-1]
			m.Returns = m.Returns[:len(m.Returns)-1]
		}
		break

	// FUNCTION_FOLLOWS jump_count
	// Adds a pointer to the constant list for the next line.
	// Jumps the given number of lines (past the function).
	case FUNCTION_FOLLOWS:
		fmt.Printf("[%v] FUNCTION_FOLLOWS %v\n", m.InstructionCounter, instruction[1])
		m.Constant = append(
			m.Constant,
			vmNumber{Type: TYPE_INT, Integer: int(m.InstructionCounter + 1)},
		)

		m.InstructionCounter += uint(instruction[1]) + 1
		break

	// CALL constant_ref
	// Jump to the instruction given in the constant.
	// Push the next instruction onto the return list.
	case CALL:
		fmt.Printf("[%v] CALL %v\n", m.InstructionCounter, instruction[1])
		jump := uint(m.Constant[instruction[1]].(vmNumber).AsInt())

		m.Returns = append(
			m.Returns,
			m.InstructionCounter+1,
		)

		m.InstructionCounter = jump
		break

	// JUMP_IF_FALSE jump_count
	// If the top constant is false, jump the number of given lines.
	case JUMP_IF_FALSE:
		fmt.Printf("[%v] JUMP_IF_FALSE %v\n", m.InstructionCounter, instruction[1])

		if m.Constant[len(m.Constant)-1].IsFalse() {
			m.InstructionCounter += uint(instruction[1])
		} else {
			m.InstructionCounter += 1
		}

		break

	// COPY constant_ref
	// Bring the constant given by constant_ref to the top of the list.
	case COPY:
		fmt.Printf("[%v] COPY %v\n", m.InstructionCounter, instruction[1])

		m.Constant = append(
			m.Constant,
			m.Constant[instruction[1]],
		)

		m.InstructionCounter += 1
		break
	}

}

func (m *Machine) Result() string {
	return fmt.Sprintf("%v", m.Constant[len(m.Constant)-1])
}

func (m *Machine) String() string {
	return fmt.Sprintf("VM@%v\nC=%v\nR=%v", m.InstructionCounter, m.Constant, m.Returns)
}
