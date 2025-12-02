package main

import (
	"fmt"
	"monkey/execution"
	"os"
	"os/user"
)

const EXT = ".mky"

func runMultipleFiles(directory string) {
	dir, err := os.ReadDir(directory)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, file := range dir {
		name := file.Name()
		if name[0] == '.' {
			continue
		}

		if file.IsDir() {
			runMultipleFiles(name)
			continue
		}

		if len(name) > len(EXT) && name[len(name)-4:] == EXT {
			execution.RunCode(os.Stdin, os.Stdout, directory+"/"+name)
		}
	}
}

func main() {
	if len(os.Args) > 1 {
		filepath := os.Args[1]

		stat, err := os.Stat(filepath)
		if err != nil {
			fmt.Println(err)
			return
		}

		if stat.IsDir() {
			runMultipleFiles(stat.Name())
		} else {
			execution.RunCode(os.Stdin, os.Stdout, filepath)
		}
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
