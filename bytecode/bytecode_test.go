package bytecode

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMake(t *testing.T) {
	tests := []struct {
		op       Opcode
		operands []int
		expected []byte
	}{
		{
			OpConstant,
			[]int{65534},
			[]byte{byte(OpConstant), 255, 254},
		},
		{
			OpAdd,
			[]int{},
			[]byte{byte(OpAdd)},
		},
		{
			OpGetLocal,
			[]int{255},
			[]byte{byte(OpGetLocal), 255},
		},
	}

	for _, test := range tests {
		instruction := Make(test.op, test.operands...)

		assert.Equal(t, len(instruction), len(test.expected), "Expected length of instruction")

		for i, _ := range test.expected {
			assert.Equal(t, instruction[i], test.expected[i], "Wrong byte in instruction")
		}
	}
}
