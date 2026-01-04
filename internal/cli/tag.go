package cli

import (
	"errors"
)

type TagCommand struct {
	*Cli
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
		noteData, err := tagCmd.Cli.GetStore().Load(tagArgs[0], tagCmd.Cli.GetUsername(), tagCmd.Cli.GetKeyProvider())
		if err != nil {
			return err
		}
		if err := noteData.AddTag(tagArgs[1]); err != nil {
			return err
		}
		if err := tagCmd.Cli.GetStore().Save(noteData, tagCmd.Cli.GetUsername(), tagCmd.Cli.GetKeyProvider()); err != nil {
			return err
		}
	case "delete","remove":
		noteData, err := tagCmd.Cli.GetStore().Load(tagArgs[0], tagCmd.Cli.GetUsername(), tagCmd.Cli.GetKeyProvider())
		if err != nil {
			return err
		}
		if err := noteData.RemoveTag(tagArgs[1]); err != nil {
			return err
		}
		if err := tagCmd.Cli.GetStore().Save(noteData, tagCmd.Cli.GetUsername(), tagCmd.Cli.GetKeyProvider()); err != nil {
			return err
		}
	case "help":
		tagCmd.Help()
	}
	return nil
}
