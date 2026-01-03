package cli

import (
	"errors"
	"fmt"
)

type SearchCommand struct {
	*Cli
}

func (searchCmd *SearchCommand) Name() string {
	return "search"
}

func (searchCmd *SearchCommand) Description() string {
	return "Search notes by keywords or tags"
}

func (searchCmd *SearchCommand) Run(args []string) error {
	if len(args) < 1 {
		searchCmd.Help()
		return errors.New("missing arguments")
	}
	cmd := args[0]
	linkArgs := args[1:]
	if len(linkArgs) < 1 {
		return errors.New("missing operand")
	}
	switch cmd {
	case "keyword", "tag":
		results, err := searchCmd.store.Search(cmd, linkArgs, searchCmd.username, searchCmd.keyProvider)
		if err != nil {
			return err
		}
		fmt.Printf("%s(s) found in %v", cmd, results)
	case "help":
		searchCmd.Help()
	}
	return nil
}
