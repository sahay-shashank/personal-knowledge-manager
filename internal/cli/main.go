package cli

import (
	"flag"
	"fmt"
	"os"

	"github.com/sahay-shashank/personal-knownledge-manager/internal/note"
)

func NewCli() {
	flag.Parse()

	commandArgs := flag.Args()
	if len(commandArgs) < 1 {
		globalHelp()
		os.Exit(1)
	}

	cli := Cli{
		store: note.InitStore("."),
	}

	cmdName := commandArgs[0]
	args := commandArgs[1:]

	commands := []Command{
		&NoteCommand{Cli: &cli},
	}
	for _, cmd := range commands {
		if cmd.Name() == cmdName {
			if err := cmd.Run(args); err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v", err)
				os.Exit(1)
			}
			return
		}
	}
	fmt.Fprintf(os.Stderr, "Unknown command: %s\n\n", cmdName)
	globalHelp()
	os.Exit(1)
}

func globalHelp() {
	fmt.Println(`
pkm: A Personal Knowledge Manager based on Zettelkasten

Usage:
  pkm <command> [arguments] [flags]

Available Command:
  note			Note function
  help			Show help information

Use "pkm help <command> for detailed information regarding a specific command.
	`)

}
