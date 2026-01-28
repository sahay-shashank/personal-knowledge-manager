package cli

import (
	"errors"
	"fmt"

	"github.com/sahay-shashank/personal-knowledge-manager/internal/crypt"
)

type UserCommand struct {
	CLI *Cli
}

func (userCmd *UserCommand) Name() string {
	return "user"
}

func (userCmd *UserCommand) Description() string {
	return "Create users, change passwords, manage accounts"
}

func (userCmd *UserCommand) Run(args []string) error {
	if len(args) < 1 {
		userCmd.Help()
		return errors.New("missing subcommand")
	}

	subCmd := args[0]
	subArgs := args[1:]

	switch subCmd {
	case "init":
		if len(subArgs) < 1 {
			return errors.New("usage: user init <username>")
		}
		return userCmd.initUser(subArgs[0])

	case "passwd":
		if len(subArgs) < 1 {
			return errors.New("usage: user passwd <username>")
		}
		return userCmd.changePassword(subArgs[0])

	case "export":
		if len(subArgs) < 1 {
			return errors.New("usage: user export <username>")
		}
		return userCmd.exportUser(subArgs[0])

	case "import":
		return userCmd.importUser()

	default:
		return fmt.Errorf("unknown subcommand: %s", subCmd)
	}
}

func (userCmd *UserCommand) initUser(username string) error {
	password, err := crypt.PromptPasswordConfirm(fmt.Sprintf("Enter password for %q: ", username))
	if err != nil {
		return err
	}

	if err := crypt.InitUser(userCmd.CLI.store.StoreLocation, username, password); err != nil {
		return err
	}

	fmt.Printf("✓ User %q initialized\n", username)
	return nil
}

func (userCmd *UserCommand) changePassword(username string) error {
	oldPassword, err := crypt.PromptPassword("Enter current password: ")
	if err != nil {
		return err
	}

	newPassword, err := crypt.PromptPasswordConfirm("Enter new password: ")
	if err != nil {
		return err
	}

	if err := crypt.ChangePassword(userCmd.CLI.store.StoreLocation, username, oldPassword, newPassword); err != nil {
		return err
	}

	fmt.Printf("✓ Password changed for %q\n", username)
	return nil
}

func (userCmd *UserCommand) importUser() error {
	userDataString, err := tempEditor(nil)
	if err != nil {
		return err
	}

	if err := crypt.ImportUser(userCmd.CLI.store.StoreLocation, userDataString); err != nil {
		return err
	}

	return nil
}

func (userCmd *UserCommand) exportUser(username string) error {
	password, err := crypt.PromptPassword(fmt.Sprintf("Enter password for %q: ", username))
	if err != nil {
		return err
	}

	if err := crypt.ExportUser(userCmd.CLI.store.StoreLocation, username, password); err != nil {
		return err
	}

	return nil
}
