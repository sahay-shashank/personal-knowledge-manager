package cli

import (
	"errors"
)

type TagCommand struct {
	*Cli
}

func (tagCmd *TagCommand) Name() string        { return "tag" }
func (tagCmd *TagCommand) Description() string { return "tag commands for note" }
func (tagCmd *TagCommand) Help() string        { return "Help section for tagging" }
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
		noteData, err := tagCmd.store.Load(tagArgs[0])
		if err != nil {
			return err
		}
		if err := noteData.AddTag(tagArgs[1]); err != nil {
			return err
		}
		if err := tagCmd.store.Save(noteData); err != nil {
			return err
		}
	case "delete":
		noteData, err := tagCmd.store.Load(tagArgs[0])
		if err != nil {
			return err
		}
		if err := noteData.RemoveTag(tagArgs[1]); err != nil {
			return err
		}
		if err := tagCmd.store.Save(noteData); err != nil {
			return err
		}
	case "help":
		tagCmd.Help()
	}
	return nil
}
