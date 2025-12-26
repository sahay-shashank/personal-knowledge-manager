package cli

import (
	"errors"

	"github.com/sahay-shashank/personal-knowledge-manager/internal/crypt"
)

type LinkCommand struct {
	*Cli
	username    string
	keyProvider *crypt.KeyProvider
}

func (linkCmd *LinkCommand) Name() string {
	return "link"
}

func (linkCmd *LinkCommand) Description() string {
	return "Create and manage links between notes for knowledge discovery"
}

func (linkCmd *LinkCommand) Run(args []string) error {
	if len(args) < 1 {
		linkCmd.Help()
		return errors.New("missing arguments")
	}
	cmd := args[0]
	linkArgs := args[1:]
	if len(linkArgs) < 2 {
		return errors.New("missing operand")
	}
	switch cmd {
	case "add":
		noteData1, err := linkCmd.store.Load(linkArgs[0], linkCmd.username, linkCmd.keyProvider)
		if err != nil {
			return err
		}
		if err := noteData1.AddLink(linkArgs[1]); err != nil {
			return err
		}
		if err := linkCmd.store.Save(noteData1, linkCmd.username, linkCmd.keyProvider); err != nil {
			return err
		}
		// back linking
		noteData2, err := linkCmd.store.Load(linkArgs[1], linkCmd.username, linkCmd.keyProvider)
		if err != nil {
			return err
		}
		if err := noteData2.AddLink(linkArgs[0]); err != nil {
			return err
		}
		if err := linkCmd.store.Save(noteData2, linkCmd.username, linkCmd.keyProvider); err != nil {
			return err
		}
	case "delete":
		noteData1, err := linkCmd.store.Load(linkArgs[0], linkCmd.username, linkCmd.keyProvider)
		if err != nil {
			return err
		}
		if err := noteData1.RemoveLink(linkArgs[1]); err != nil {
			return err
		}
		if err := linkCmd.store.Save(noteData1, linkCmd.username, linkCmd.keyProvider); err != nil {
			return err
		}
		// back linking
		noteData2, err := linkCmd.store.Load(linkArgs[1], linkCmd.username, linkCmd.keyProvider)
		if err != nil {
			return err
		}
		if err := noteData2.RemoveLink(linkArgs[0]); err != nil {
			return err
		}
		if err := linkCmd.store.Save(noteData2, linkCmd.username, linkCmd.keyProvider); err != nil {
			return err
		}
	case "help":
		linkCmd.Help()
	}
	return nil
}
