package vm

import (
	"go_interpreter/bytecode"
	"go_interpreter/object"
)

// Holds data relevant to execution
type Frame struct {
	fn *object.CompiledFunction // Compiled function referenced by frame
	ip int                      // Instruction pointer to the compiled function
}

func BuildFrame(fn *object.CompiledFunction) *Frame {
	return &Frame{fn, -1}
}

func (f *Frame) Instructions() bytecode.Instructions {
	return f.fn.Instructions
}
