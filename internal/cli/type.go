package cli

import "github.com/sahay-shashank/personal-knownledge-manager/internal/note"

type Command interface {
	Name() string
	Description() string
	Run(args []string) error
	Help() string
}

type Cli struct {
	store *note.Store
}
