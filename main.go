package main

import (
	"fmt"
	"go_interpreter/repl"
	"os"
	"os/user"
)

func main() {
	// Get user
	user, err := user.Current()
	if err != nil {
		panic(err)
	}

	// Print welcome prompt
	fmt.Printf("Welcome to the Monkey programming language, %s!\n", user.Username)
	fmt.Printf("Feel free to type in commands\n")

	// Start loop
	repl.StartLoop(os.Stdin, os.Stdout)
}
