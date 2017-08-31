package main

import (
	"GoClang/repl"
	"fmt"
	"os"
	"os/user"
)

func main() {
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Hello, %s! This is The GoClang Programming Language.\n", usr.Username)
	fmt.Printf("Feel free to type in command!\n")
	repl.Start(os.Stdin, os.Stdout)
}
