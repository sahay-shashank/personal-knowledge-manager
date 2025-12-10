package cli

import (
	"errors"
)

type LinkCommand struct {
	*Cli
}

func (linkCmd *LinkCommand) Name() string        { return "link" }
func (linkCmd *LinkCommand) Description() string { return "Link commands for note" }
func (linkCmd *LinkCommand) Help() string        { return "Help section for linking" }
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
		noteData1, err := linkCmd.store.Load(linkArgs[0])
		if err != nil {
			return err
		}
		if err := noteData1.AddLink(linkArgs[1]); err != nil {
			return err
		}
		if err := linkCmd.store.Save(noteData1); err != nil {
			return err
		}
		// back linking
		noteData2, err := linkCmd.store.Load(linkArgs[1])
		if err != nil {
			return err
		}
		if err := noteData2.AddLink(linkArgs[0]); err != nil {
			return err
		}
		if err := linkCmd.store.Save(noteData2); err != nil {
			return err
		}
	case "delete":
		noteData1, err := linkCmd.store.Load(linkArgs[0])
		if err != nil {
			return err
		}
		if err := noteData1.RemoveLink(linkArgs[1]); err != nil {
			return err
		}
		if err := linkCmd.store.Save(noteData1); err != nil {
			return err
		}
		// back linking
		noteData2, err := linkCmd.store.Load(linkArgs[1])
		if err != nil {
			return err
		}
		if err := noteData2.RemoveLink(linkArgs[0]); err != nil {
			return err
		}
		if err := linkCmd.store.Save(noteData2); err != nil {
			return err
		}
	case "help":
		linkCmd.Help()
	}
	return nil
}
