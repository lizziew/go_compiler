package bytecode

import (
	"encoding/binary"
	"fmt"
)

type Instructions []byte
type Opcode byte

const (
	OpConstant      Opcode = iota // 1 operand: previous assigned number to constant
	OpAdd                         // 0 operands
	OpPop                         // 0 operands
	OpSub                         // 0 operands
	OpMul                         // 0 operands
	OpDiv                         // 0 operands
	OpTrue                        // 0 operands
	OpFalse                       // 0 operands
	OpEqual                       // 0 operands
	OpNotEqual                    // 0 operands
	OpGreater                     // 0 operands
	OpMinus                       // 0 operands
	OpBang                        // 0 operands
	OpJumpNotTruthy               // 1 operand: jump offset if stack top is false, not null
	OpJump                        // 1 operand: jump offset)
	OpNull                        // 0 operands
	OpGetGlobal                   // 1 operand: unique index of global binding
	OpSetGlobal                   // 1 operand: unique index of global binding
	OpArray                       // 1 operand: number of elements
	OpHash                        // 1 operand: number of key + value elements
	OpIndex                       // 0 operands
	OpCall                        // 0 operands
	OpReturnValue                 // 0 operands: return value at top of stack
	OpReturnNothing               // 0 operands: return from current function (no value)
	OpSetLocal                    // 1 operand: unique index of local binding
	OpGetLocal                    // 1 operand: unique index of local binding
)

type Definition struct {
	Name          string // readability
	OperandWidths []int  // number of bytes each operand takes up
}

var definitions = map[Opcode]*Definition{
	OpConstant:      {"OpConstant", []int{2}},
	OpAdd:           {"OpAdd", []int{}},
	OpPop:           {"OpPop", []int{}},
	OpSub:           {"OpSub", []int{}},
	OpMul:           {"OpMul", []int{}},
	OpDiv:           {"OpDiv", []int{}},
	OpTrue:          {"OpTrue", []int{}},
	OpFalse:         {"OpFalse", []int{}},
	OpEqual:         {"OpEqual", []int{}},
	OpNotEqual:      {"OpNotEqual", []int{}},
	OpGreater:       {"OpGreater", []int{}},
	OpMinus:         {"OpMinus", []int{}},
	OpBang:          {"OpBang", []int{}},
	OpJumpNotTruthy: {"OpJumpNotTruthy", []int{2}},
	OpJump:          {"OpJump", []int{2}},
	OpNull:          {"OpNull", []int{}},
	OpGetGlobal:     {"OpGetGlobal", []int{2}},
	OpSetGlobal:     {"OpSetGlobal", []int{2}},
	OpArray:         {"OpArray", []int{2}},
	OpHash:          {"OpHash", []int{2}},
	OpIndex:         {"OpIndex", []int{}},
	OpCall:          {"OpCall", []int{}},
	OpReturnValue:   {"OpReturnValue", []int{}},
	OpReturnNothing: {"OpReturnNothing", []int{}},
	OpGetLocal:      {"OpGetLocal", []int{1}},
	OpSetLocal:      {"OpSetLocal", []int{1}},
}

// Make instruction from op and operands (Big Endian)
func Make(op Opcode, operands ...int) []byte {
	definition, ok := definitions[op]
	if !ok {
		return []byte{}
	}

	numBytes := 1
	for _, width := range definition.OperandWidths {
		numBytes += width
	}

	instruction := make([]byte, numBytes)

	// Add op
	instruction[0] = byte(op)

	// Add operands
	offset := 1
	for i, o := range operands {
		width := definition.OperandWidths[i]

		switch width {
		case 1:
			instruction[offset] = byte(o)
		case 2:
			binary.BigEndian.PutUint16(instruction[offset:], uint16(o))
		}

		offset += width
	}

	return instruction
}

func ReadUint16(i Instructions) uint16 {
	return binary.BigEndian.Uint16(i)
}

func ReadUint8(i Instructions) uint8 {
	return uint8(i[0])
}

// For debugging
func Lookup(op byte) (*Definition, error) {
	definition, ok := definitions[Opcode(op)]
	if !ok {
		return nil, fmt.Errorf("Opcode %d undefined", op)
	}

	return definition, nil
}
