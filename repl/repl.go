package repl

import (
	"bufio"
	"fmt"
	"go_interpreter/lexer"
	"go_interpreter/token"
	"io"
)

const PROMPT = ">> "

func StartLoop(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Printf(PROMPT)

		// Get user input
		scanned := scanner.Scan()

		// Stop when newline is encountered
		if !scanned {
			return
		}

		// Feed input into lexer
		l := lexer.BuildLexer(scanner.Text())
		for {
			t := l.NextToken()

			if t.Type == token.EOF {
				break
			}

			fmt.Printf("%+v\n", t)
		}
	}
}
