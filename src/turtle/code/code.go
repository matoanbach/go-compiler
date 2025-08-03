package code

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type Instructions []byte

type Opcode byte

const (
	OpConstant Opcode = iota
	OpAdd
	OpPop
	OpSub
	OpMul
	OpDiv
	OpTrue
	OpFalse
	OpEqual
	OpNotEqual
	OpGreaterThan
	OpMinus
	OpBang
	OpJumpNotTruthy
	OpJump // jump no matter what
	OpNull
	OpSetGlobal
	OpGetGlobal
	OpArray
	OpHash
	OpIndex
	OpCall
	OpReturnValue
	OpReturn
	OpSetLocal
	OpGetLocal
	OpGetBuiltin
	OpClosure
)

type Definition struct {
	Name          string
	OperandWidths []int
}

var definitions = map[Opcode]*Definition{
	OpConstant:      {"OpConstant", []int{2}},
	OpAdd:           {"OpAdd", []int{}},
	OpPop:           {"OpPop", []int{}},
	OpSub:           {"ObSub", []int{}},
	OpMul:           {"OpMul", []int{}},
	OpDiv:           {"OpDiv", []int{}},
	OpTrue:          {"OpTrue", []int{}},
	OpFalse:         {"OpFalse", []int{}},
	OpEqual:         {"OpEqual", []int{}},
	OpNotEqual:      {"OpNotEqual", []int{}},
	OpGreaterThan:   {"OpGreaterThan", []int{}},
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
	OpCall:          {"OpCall", []int{1}},
	OpReturn:        {"OpReturn", []int{}},
	OpReturnValue:   {"OpReturnValue", []int{}},
	OpSetLocal:      {"OpSetLocal", []int{1}},
	OpGetLocal:      {"OpGetLocal", []int{1}},
	OpGetBuiltin:    {"OpGetBuiltin", []int{1}},
	OpClosure:       {"OpClosure", []int{2, 1}},
}

func Lookup(op byte) (*Definition, error) {
	def, ok := definitions[Opcode(op)]
	if !ok {
		return nil, fmt.Errorf("opcode %d underfined", op)
	}

	return def, nil
}

func Make(op Opcode, operands ...int) []byte {
	def, ok := definitions[op]
	if !ok {
		return []byte{}
	}

	// count the instructionLen
	instructionLen := 1
	for _, w := range def.OperandWidths {
		instructionLen += w
	}

	// loop through each operand and put it in the slice "instruction"
	offset := 1
	instruction := make([]byte, instructionLen)
	instruction[0] = byte(op)
	for i, operand := range operands {
		width := def.OperandWidths[i]
		switch width {
		case 2:
			binary.BigEndian.PutUint16(instruction[offset:], uint16(operand))
		case 1:
			instruction[offset] = byte(operand)
		}
		offset += width
	}

	return instruction
}

func ReadOperands(def *Definition, instruction Instructions) ([]int, int) {
	operands := make([]int, len(def.OperandWidths))
	offset := 0

	for i, width := range def.OperandWidths {
		switch width {
		case 1:
			operands[i] = int(ReadUint8(instruction[offset:]))
		case 2:
			operands[i] = int(ReadUint16(instruction[offset:]))
		}
		offset += width
	}
	return operands, offset
}

func ReadUint8(instruction Instructions) uint8 {
	return byte(instruction[0])
}
func ReadUint16(instruction Instructions) uint16 {
	return binary.BigEndian.Uint16(instruction)
}

func (instruction Instructions) String() string {
	var out bytes.Buffer
	// if errr, fmt.Printf(&out, "ERROR: %s\n", err)
	i := 0
	for i < len(instruction) {
		def, err := Lookup(instruction[i])
		if err != nil {
			fmt.Fprintf(&out, "ERROR: %s\n", err)
		}

		operands, read := ReadOperands(def, instruction[i+1:])
		fmt.Fprintf(&out, "%04d %s\n", i, instruction.fmtInstruction(def, operands))

		i += 1 + read
	}

	return out.String()
}

func (ins Instructions) fmtInstruction(def *Definition, operands []int) string {
	operandCount := len(operands)
	if operandCount != len(def.OperandWidths) {
		return fmt.Sprintf("ERROR: operand len %d does not match defined %d\n", len(operands), operandCount)
	}

	switch operandCount {
	case 0:
		return fmt.Sprintf("%s", def.Name) // or return def.Name
	case 1:
		return fmt.Sprintf("%s %d", def.Name, operands[0])
	}

	return fmt.Sprintf("ERROR: unhandled operandCount for %s\n", def.Name)
}
