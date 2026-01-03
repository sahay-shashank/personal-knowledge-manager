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

// GetStore returns the note store
func (c *Cli) GetStore() *note.Store {
	return c.store
}

// GetUsername returns the username
func (c *Cli) GetUsername() string {
	return c.username
}

// GetKeyProvider returns the key provider
func (c *Cli) GetKeyProvider() *crypt.KeyProvider {
	return c.keyProvider
}

// SetStore sets the note store (for testing)
func (c *Cli) SetStore(s *note.Store) {
	c.store = s
}

// SetUsername sets the username (for testing)
func (c *Cli) SetUsername(u string) {
	c.username = u
}

// SetKeyProvider sets the key provider (for testing)
func (c *Cli) SetKeyProvider(kp *crypt.KeyProvider) {
	c.keyProvider = kp
}
