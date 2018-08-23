package repl

import (
	"bufio"
	"fmt"
	"go_interpreter/compiler"
	// "go_interpreter/evaluator"
	"go_interpreter/lexer"
	"go_interpreter/object"
	"go_interpreter/parser"
	"go_interpreter/vm"
	"io"
)

const PROMPT = ">> "

func StartLoop(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	// env := object.BuildEnvironment()

	constants := []object.Object{}
	globals := make([]object.Object, vm.GlobalCapacity)
	symbolTable := compiler.BuildSymbolTable()

	for {
		fmt.Printf(PROMPT)

		// Get user input
		scanned := scanner.Scan()

		// Stop when newline is encountered
		if !scanned {
			return
		}

		// Lexer
		l := lexer.BuildLexer(scanner.Text())

		// Parser
		p := parser.BuildParser(l)
		prog := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		// Evaluator
		/*result := evaluator.Eval(prog, env)
		if result != nil {
			io.WriteString(out, result.Inspect())
			io.WriteString(out, "\n")
		}*/

		// Compiler
		c := compiler.BuildStatefulCompiler(symbolTable, constants)
		err := c.Compile(prog)
		if err != nil {
			fmt.Fprintf(out, "Compile-time error: %s\n", err)
		}

		// VM
		bytecode := c.Bytecode()
		constants = bytecode.Constants
		machine := vm.BuildStatefulVM(bytecode, globals)
		err = machine.Run()
		if err != nil {
			fmt.Fprintf(out, "Run-time error: %s\n", err)
		}
		lastPopped := machine.LastPopped()
		io.WriteString(out, lastPopped.Inspect())
		io.WriteString(out, "\n")
	}
}

func printParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
