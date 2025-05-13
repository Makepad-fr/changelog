package internal

import (
	"fmt"
	"log"
	"os"

	"github.com/Makepad-fr/changelog/parser"
)

func VerifyCommand() {
	args := os.Args[2:]

	file := "CHANGELOG.md"
	if len(args) > 0 && args[0] != "" {
		file = args[0]
	}

	cl, err := parser.LoadChangelogFromFile(file)
	if err != nil {
		log.Fatalf("Can not load changelog from file %s: %v", file, err)
	}
	err = (*cl).Verify()
	if err != nil {
		log.Fatalf("Changelog is not valid: %v", err)
	}
	fmt.Printf("%s is a valid.\n", file)
}
