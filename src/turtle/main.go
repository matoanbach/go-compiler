package main

import (
	"os"
	"turtle/repl"
)

func main() {
	repl.Start(os.Stdin, os.Stdout)
}
