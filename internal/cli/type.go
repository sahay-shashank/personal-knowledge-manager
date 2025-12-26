package cli

import (
	"github.com/sahay-shashank/personal-knowledge-manager/internal/crypt"
	"github.com/sahay-shashank/personal-knowledge-manager/internal/note"
)

type Command interface {
	Name() string
	Description() string
	Run(args []string) error
	Help() string
}

type Cli struct {
	store       *note.Store
	username    string
	keyProvider *crypt.KeyProvider
}
