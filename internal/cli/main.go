package cli

import (
	"flag"
	"fmt"
	"os"

	"github.com/sahay-shashank/personal-knowledge-manager/internal/crypt"
	"github.com/sahay-shashank/personal-knowledge-manager/internal/note"
)

func NewCli() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v", err)
	}

	storeDir := flag.String("storeDirectory", homeDir+"/.pkm/", "Note Storage Directory")
	username := flag.String("user", "", "Username (required)")

	flag.Parse()

	commandArgs := flag.Args()
	if len(commandArgs) < 1 {
		globalHelp()
		os.Exit(1)
	}

	cmdName := commandArgs[0]
	args := commandArgs[1:]

	// Special handling for 'user' commands (no password needed)
	if cmdName == "user" {
		cli := &Cli{
			store:    note.InitStore(*storeDir),
			username: "",
		}
		cmd := &UserCommand{CLI: cli}
		if err := cmd.Run(args); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		return
	}

	if *username == "" {
		fmt.Fprintln(os.Stderr, "Error: --user flag is required")
		os.Exit(1)
	}

	// Prompt for password (NEW)
	password, err := crypt.PromptPassword("Password: ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Initialize crypto with hardcoded DEK (for testing)
	keyProvider, err := crypt.NewKeyProvider(*storeDir, *username, password)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Crypto init error: %v\n", err)
		os.Exit(1)
	}

	cli := Cli{
		store:       note.InitStore(*storeDir),
		username:    *username,
		keyProvider: keyProvider,
	}

	commands := []Command{
		&NoteCommand{Cli: &cli, username: *username, keyProvider: keyProvider},
		&LinkCommand{Cli: &cli},
		&TagCommand{Cli: &cli},
		&SearchCommand{Cli: &cli},
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
