package main

import (
	"fmt"
	"monkey/execution"
	"os"
	"os/user"
)

func main() {
	if len(os.Args) > 1 {
		execution.RunCode(os.Stdin, os.Stdout, os.Args[1])
		return
	}

	user, err := user.Current()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Hi %s! This is the Monkey Programming Language\n", user.Username)
	fmt.Printf("Feel free to type in commands\n")

	execution.StartREPL(os.Stdin, os.Stdout)
}
