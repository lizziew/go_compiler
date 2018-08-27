package main

import (
	"flag"
	"fmt"
	"go_interpreter/compiler"
	"go_interpreter/evaluator"
	"go_interpreter/lexer"
	"go_interpreter/object"
	"go_interpreter/parser"
	"go_interpreter/vm"
	"time"
)

var engine = flag.String("engine", "vm", "use 'vm' or 'eval'")

var input = `
let fibonacci = fn(x) {
  if (x == 0) {
    0
  } else {
    if (x == 1) {
      return 1;
    } else {
      fibonacci(x - 1) + fibonacci(x - 2);
    }
  }
};
fibonacci(35);
`

func main() {
	flag.Parse()

	var duration time.Duration
	var result object.Object

	l := lexer.BuildLexer(input)
	p := parser.BuildParser(l)
	program := p.ParseProgram()

	if *engine == "vm" {
		comp := compiler.BuildCompiler()
		err := comp.Compile(program)
		if err != nil {
			fmt.Printf("compiler error: %s", err)
			return
		}

		machine := vm.BuildVM(comp.Bytecode())

		start := time.Now()

		err = machine.Run()
		if err != nil {
			fmt.Printf("vm error: %s", err)
			return
		}

		duration = time.Since(start)
		result = machine.LastPopped()
	} else {
		env := object.BuildEnvironment()
		start := time.Now()
		result = evaluator.Eval(program, env)
		duration = time.Since(start)
	}

	fmt.Printf(
		"engine=%s, result=%s, duration=%s\n",
		*engine,
		result.Inspect(),
		duration)
}
