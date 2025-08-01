package code

import (
	"fmt"
	"testing"
)

func TestMake(t *testing.T) {
	tests := []struct {
		op       Opcode
		operants []int
		expected []byte
	}{
		{OpConstant, []int{65534}, []byte{byte(OpConstant), 255, 254}},
		{OpAdd, []int{}, []byte{byte(OpAdd)}},
		{OpGetLocal, []int{255}, []byte{byte(OpGetLocal), 255}},
	}
	for _, tt := range tests {
		instruction := Make(tt.op, tt.operants...)

		if len(instruction) != len(tt.expected) {
			t.Errorf("instruction has wrong length. expect=%d, got=%d", len(instruction), len(tt.expected))
		}

		for i, bb := range instruction {
			if instruction[i] != tt.expected[i] {
				t.Errorf("instruction[%d] has the wrong value. expect=%d, got=%d", i, bb, instruction[i])
			}
		}
	}
}

func TestInstructionsString(t *testing.T) {
	instructions := []Instructions{
		Make(OpAdd),
		Make(OpGetLocal, 1),
		Make(OpConstant, 2),
		Make(OpConstant, 65535),
	}

	expected := `0000 OpAdd
0001 OpGetLocal 1
0003 OpConstant 2
0006 OpConstant 65535
`
	concatted := Instructions{}
	for _, ins := range instructions {
		concatted = append(concatted, ins...)
	}

	if concatted.String() != expected {
		t.Errorf("instructions wrongly formatted. \nexpected=%q\ngot=%q", expected, concatted.String())
	}
	fmt.Printf("got=%q\n", concatted.String())
}

func TestReadOperants(t *testing.T) {
	tests := []struct {
		op        Opcode
		operants  []int
		bytesRead int
	}{
		{OpConstant, []int{65535}, 2},
	}

	for _, tt := range tests {
		instruction := Make(tt.op, tt.operants...)
		def, err := Lookup(byte(tt.op))
		if err != nil {
			t.Errorf("Opcode not found. got=%d", tt.op)
		}
		operandsRead, n := ReadOperands(def, instruction[1:])

		if n != tt.bytesRead {
			t.Fatalf("n wrong. expected=%d, got=%d", tt.bytesRead, n)
		}

		for i, want := range tt.operants {
			if operandsRead[i] != want {
				t.Errorf("operand wrong. expected=%d, got=%d", want, operandsRead[i])
			}
		}
	}
}
