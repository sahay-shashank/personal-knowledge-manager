package cli

import (
	"errors"
	"strings"
)

type LinkCommand struct {
	*Cli
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
		noteData1, err := linkCmd.Cli.GetStore().Load(linkArgs[0], linkCmd.Cli.GetUsername(), linkCmd.Cli.GetKeyProvider())
		if err != nil {
			return err
		}
		if err := noteData1.AddLink(linkArgs[1]); err != nil {
			return err
		}
		if err := linkCmd.Cli.GetStore().Save(noteData1, linkCmd.Cli.GetUsername(), linkCmd.Cli.GetKeyProvider()); err != nil {
			return err
		}
		// back linking
		noteData2, err := linkCmd.Cli.GetStore().Load(linkArgs[1], linkCmd.Cli.GetUsername(), linkCmd.Cli.GetKeyProvider())
		if err != nil {
			return err
		}
		if err := noteData2.AddLink(linkArgs[0]); err != nil {
			return err
		}
		if err := linkCmd.Cli.GetStore().Save(noteData2, linkCmd.Cli.GetUsername(), linkCmd.Cli.GetKeyProvider()); err != nil {
			return err
		}
	case "delete","remove":
		noteData1, err := linkCmd.Cli.GetStore().Load(linkArgs[0], linkCmd.Cli.GetUsername(), linkCmd.Cli.GetKeyProvider())
		if err != nil {
			return err
		}
		if err := noteData1.RemoveLink(linkArgs[1]); err != nil {
			return err
		}
		if err := linkCmd.Cli.GetStore().Save(noteData1, linkCmd.Cli.GetUsername(), linkCmd.Cli.GetKeyProvider()); err != nil {
			return err
		}
		// back linking
		noteData2, err := linkCmd.Cli.GetStore().Load(linkArgs[1], linkCmd.Cli.GetUsername(), linkCmd.Cli.GetKeyProvider())
		if err != nil {
			return err
		}
		if err := noteData2.RemoveLink(linkArgs[0]); err != nil {
			if strings.Contains(err.Error(), "link not found") {
				// ignore if back-link didn't exist
			} else {
				return err
			}
		}
		if err := linkCmd.Cli.GetStore().Save(noteData2, linkCmd.Cli.GetUsername(), linkCmd.Cli.GetKeyProvider()); err != nil {
			return err
		}
	case "help":
		linkCmd.Help()
	}
	return nil
}
