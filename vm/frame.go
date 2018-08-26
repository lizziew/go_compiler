package vm

import (
	"go_interpreter/bytecode"
	"go_interpreter/object"
)

// Holds data relevant to execution
type Frame struct {
	fn          *object.CompiledFunction // Compiled function referenced by frame
	ip          int                      // Instruction pointer to the compiled function
	basePointer int                      // Bottom of stack of current call frame
}

func BuildFrame(fn *object.CompiledFunction, basePointer int) *Frame {
	return &Frame{fn, -1, basePointer}
}

func (f *Frame) Instructions() bytecode.Instructions {
	return f.fn.Instructions
}
