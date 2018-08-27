package main

import (
	"flag"
	"fmt"
	"go_interpreter/repl"
	"os"
	"os/user"
)

func main() {
	// Interpreter or compiler
	engine := flag.String("engine", "vm", "use 'vm' or 'eval'")
	flag.Parse()

	// Get user
	user, err := user.Current()
	if err != nil {
		panic(err)
	}

	// Print welcome prompt
	fmt.Printf("Welcome to the Monkey programming language, %s!\n", user.Username)
	fmt.Printf("Feel free to type in commands. Engine = %s\n", *engine)

	// Start loop
	repl.StartLoop(engine, os.Stdin, os.Stdout)
}
