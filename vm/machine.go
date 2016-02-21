package vm

import (
	"bytes"
	"encoding/binary"
	"fmt"
	. "github.com/jonnyarnold/fn-go/bytecode"
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

	// CurrentByte tracks the current 'head' of the program.
	CurrentByte uint64

	// Running is set to true when started, and is used to
	// break out of the instruction loop when complete.
	Running bool

	// The set of return points.
	// When a RETURN instruction is met,
	// the top of this stack replaces the instruction counter.
	Returns []uint64
}

// Initializes a new Machine.
func NewMachine(bytecode []byte) Machine {
	return Machine{
		Code:        bytecode,
		Running:     false,
		CurrentByte: 0,
	}
}

// Processes all instructions in Code
// until the machines stops Running.
func (m *Machine) ProcessAll() {
	m.Running = true

	for m.Running {
		m.ExecuteNextInstruction()
	}
}

func (m *Machine) ExecuteNextInstruction() {
	fmt.Printf("[I%04x] ", m.CurrentByte)

	// The first byte is the instruction type.
	switch m.Code[m.CurrentByte] {

	// DECLARE constant
	// Adds the constant to the Constant List.
	case DECLARE:
		fmt.Print("DECLARE ")
		constType := m.Code[m.CurrentByte+1]

		switch constType {
		case TYPE_INT, TYPE_FLOAT:
			value := m.Code[m.CurrentByte+2 : m.CurrentByte+10]

			m.Constant = append(m.Constant, VMConstant(constType, value))
			m.CurrentByte += 10 // DECLARE + TYPE + VALUE
			break

		case TYPE_TRUE, TYPE_FALSE:
			m.Constant = append(m.Constant, VMConstant(constType, nil))
			m.CurrentByte += 2 // DECLARE + TYPE
			break

		case TYPE_STRING:
			valueLengthBytes := m.Code[m.CurrentByte+2 : m.CurrentByte+10]
			byteBuffer := bytes.NewReader(valueLengthBytes)
			valueLength, err := binary.ReadUvarint(byteBuffer)
			if err != nil {
				panic(err)
			}

			value := m.Code[m.CurrentByte+10 : m.CurrentByte+10+valueLength]

			m.Constant = append(m.Constant, VMConstant(constType, value))
			m.CurrentByte += valueLength + 10 // DECLARE + TYPE + LENGTH + VALUE
			break

		default:
			panic("Unknown DECLARE type.")
		}

		fmt.Printf("%v\n", m.Constant[len(m.Constant)-1])

		break

	// ADD
	// Adds the top 2 constants together and pushes them onto the Constant List.
	case ADD:
		fmt.Println("ADD")

		m.Constant = append(
			m.Constant,
			AddNumbers(m.Constant[len(m.Constant)-1].(vmNumber), m.Constant[len(m.Constant)-2].(vmNumber)))

		m.CurrentByte += 1
		break

	// SUBTRACT
	// Subtracts the top 2 constants together and pushes them onto the Constant List.
	// (If the list is [2 1], 2 - 1 will be calculated.)
	case SUBTRACT:
		fmt.Println("SUBTRACT")

		m.Constant = append(
			m.Constant,
			SubtractNumbers(m.Constant[len(m.Constant)-2].(vmNumber), m.Constant[len(m.Constant)-1].(vmNumber)))

		m.CurrentByte += 1
		break

	// MULTIPLY
	// Multiplies the top 2 constants together and pushes them onto the Constant List.
	case MULTIPLY:
		fmt.Println("MULTIPLY")

		m.Constant = append(
			m.Constant,
			MultiplyNumbers(m.Constant[len(m.Constant)-1].(vmNumber), m.Constant[len(m.Constant)-2].(vmNumber)))

		m.CurrentByte += 1
		break

	// DIVIDE
	// Divides the top 2 constants together and pushes them onto the Constant List.
	// (If the list is [2 1], 2 / 1 will be calculated.)
	case DIVIDE:
		fmt.Println("DIVIDE")

		m.Constant = append(
			m.Constant,
			DivideNumbers(m.Constant[len(m.Constant)-2].(vmNumber), m.Constant[len(m.Constant)-1].(vmNumber)))

		m.CurrentByte += 1
		break

	// AND
	// Takes the logical AND of the top 2 constants.
	case AND:
		fmt.Println("AND")
		and := m.Constant[len(m.Constant)-1].(vmBool).Value && m.Constant[len(m.Constant)-2].(vmBool).Value
		m.Constant = append(
			m.Constant,
			vmBool{Value: and})

		m.CurrentByte += 1
		break

	// OR
	// Takes the logical OR of the top 2 constants.
	case OR:
		fmt.Println("OR")
		or := m.Constant[len(m.Constant)-1].(vmBool).Value || m.Constant[len(m.Constant)-2].(vmBool).Value
		m.Constant = append(
			m.Constant,
			vmBool{Value: or})

		m.CurrentByte += 1
		break

	// NEGATE
	// Takes the logical NOT or unary - of the top constant.
	case NEGATE:
		fmt.Println("NEGATE")
		m.Constant = append(
			m.Constant,
			m.Constant[len(m.Constant)-1].Negate())

		m.CurrentByte += 1
		break

	// DUMP
	// Dumps the current machine state
	case DUMP:
		fmt.Println("DUMP")
		fmt.Println(m)
		m.CurrentByte += 1
		break

	// RETURN
	// If there is nothing on the Return list, terminate.
	// If there is a Return, pop the top return from the stack
	// and go to it.
	case RETURN:
		fmt.Print("RETURN")
		m.CurrentByte += 1

		if len(m.Returns) == 0 {
			fmt.Print("\n")
			m.Running = false
		} else {
			returnJump := m.Returns[len(m.Returns)-1]
			fmt.Printf(" -> I%04x\n", returnJump)
			m.CurrentByte = returnJump
			m.Returns = m.Returns[:len(m.Returns)-1]
		}
		break

	// FUNCTION_FOLLOWS jump_count
	// Adds a pointer to the constant list for the next line.
	// Jumps the given number of lines (past the function).
	case FUNCTION_FOLLOWS:
		jumpCountBytes := m.Code[m.CurrentByte+1 : m.CurrentByte+9]
		byteBuffer := bytes.NewReader(jumpCountBytes)
		jumpCount, err := binary.ReadUvarint(byteBuffer)
		if err != nil {
			panic(err)
		}

		fmt.Printf("FUNCTION_FOLLOWS (JUMP +%v)\n", jumpCount)

		m.Constant = append(
			m.Constant,
			vmNumber{Type: TYPE_INT, Integer: int64(m.CurrentByte + 9)},
		)

		m.CurrentByte += jumpCount + 9
		break

	// CALL constant_ref
	// Jump to the instruction given in the constant.
	// Push the next instruction onto the return list.
	case CALL:
		constIdx := m.Code[m.CurrentByte+1]
		fmt.Printf("CALL C%02x", constIdx)

		jump := uint64(m.Constant[constIdx].(vmNumber).AsInt())

		fmt.Printf(" -> I%04x\n", jump)

		m.Returns = append(
			m.Returns,
			m.CurrentByte+2,
		)

		m.CurrentByte = jump
		break

	// JUMP_IF_FALSE jump_count
	// If the top constant is false, jump the number of given bytes.
	case JUMP_IF_FALSE:
		if m.Constant[len(m.Constant)-1].IsFalse() {
			jumpCountBytes := m.Code[m.CurrentByte+1 : m.CurrentByte+9]
			byteBuffer := bytes.NewReader(jumpCountBytes)
			jumpCount, err := binary.ReadUvarint(byteBuffer)
			if err != nil {
				panic(err)
			}

			fmt.Printf("JUMP_IF_FALSE (JUMP +%v)\n", jumpCount)
			m.CurrentByte += uint64(jumpCount) + 9
		} else {
			fmt.Println("JUMP_IF_FALSE (NOOP)")
			m.CurrentByte += 9
		}

		break

	// COPY constant_ref
	// Bring the constant given by constant_ref to the top of the list.
	case COPY:
		constIdx := m.Code[m.CurrentByte+1]
		fmt.Printf("COPY C%02x\n", constIdx)

		m.Constant = append(
			m.Constant,
			m.Constant[constIdx],
		)

		m.CurrentByte += 2
		break

	default:
		panic(fmt.Sprintf("Unknown byte %02x", m.Code[m.CurrentByte]))
	}

}

func (m *Machine) Result() string {
	return fmt.Sprintf("%v", m.Constant[len(m.Constant)-1])
}

func (m *Machine) String() string {
	return fmt.Sprintf("VM@I%04x\nC=%v\nR=%v", m.CurrentByte, m.Constant, m.Returns)
}
