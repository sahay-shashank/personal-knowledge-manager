package cli

import (
	"errors"

	"github.com/sahay-shashank/personal-knowledge-manager/internal/crypt"
)

type TagCommand struct {
	*Cli
	username    string
	keyProvider *crypt.KeyProvider
}

func (tagCmd *TagCommand) Name() string {
	return "tag"
}

func (tagCmd *TagCommand) Description() string {
	return "Organize notes with tags for categorization and search"
}

func (tagCmd *TagCommand) Run(args []string) error {
	if len(args) < 1 {
		tagCmd.Help()
		return errors.New("missing arguments")
	}
	cmd := args[0]
	tagArgs := args[1:]
	if len(tagArgs) < 2 {
		return errors.New("missing operand")
	}
	switch cmd {
	case "add":
		noteData, err := tagCmd.store.Load(tagArgs[0], tagCmd.username, tagCmd.keyProvider)
		if err != nil {
			return err
		}
		if err := noteData.AddTag(tagArgs[1]); err != nil {
			return err
		}
		if err := tagCmd.store.Save(noteData, tagCmd.username, tagCmd.keyProvider); err != nil {
			return err
		}
	case "delete":
		noteData, err := tagCmd.store.Load(tagArgs[0], tagCmd.username, tagCmd.keyProvider)
		if err != nil {
			return err
		}
		if err := noteData.RemoveTag(tagArgs[1]); err != nil {
			return err
		}
		if err := tagCmd.store.Save(noteData, tagCmd.username, tagCmd.keyProvider); err != nil {
			return err
		}
	case "help":
		tagCmd.Help()
	}
	return nil
}
