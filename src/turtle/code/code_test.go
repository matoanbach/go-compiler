package code

import (
	"testing"
)

func TestMake(t *testing.T) {
	tests := []struct {
		op       Opcode
		operants []int
		expected []byte
	}{
		{OpConstant, []int{65534}, []byte{byte(OpConstant), 255, 254}},
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
