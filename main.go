package main

import (
	"os"
	"github.com/ocfl/ocfl/commands"
)

func main() {
	commands.Execute(os.Args[1:])
}
