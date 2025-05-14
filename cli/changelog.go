package main

import (
	"fmt"
	"os"

	"github.com/Makepad-fr/changelog/cli/internal"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: changelog <command> [options]")
		fmt.Println("Commands: verify, new")
		os.Exit(1)
	}

	command := os.Args[1]
	switch command {
	case "validate":
		internal.ValidateCommand()
	case "new":
		internal.NewCommand()
	default:
		fmt.Printf("Unknown command: %s\\n", command)
		os.Exit(1)
	}
}
