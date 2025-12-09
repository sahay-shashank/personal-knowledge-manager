package cli

import (
	"errors"
	"os"
	"os/exec"

	"github.com/sahay-shashank/personal-knownledge-manager/internal/note"
)

type NoteCommand struct {
	*Cli
}

func (noteCmd *NoteCommand) Name() string        { return "note" }
func (noteCmd *NoteCommand) Description() string { return "CRUD commands for note" }
func (noteCmd *NoteCommand) Help() string        { return "Help section for note" }
func (noteCmd *NoteCommand) Run(args []string) error {
	if len(args) < 1 {
		noteCmd.Help()
		return errors.New("missing arguments")
	}
	cmd := args[0]
	noteArgs := args[1:]
	switch cmd {
	case "new":
		content, err := tempEditor(nil)
		if err != nil {
			return err
		}
		noteData := note.NewNote(noteArgs[0], content)
		if err := noteCmd.store.Save(noteData); err != nil {
			return err
		}
	case "edit":
		noteData, err := noteCmd.store.Load(noteArgs[0])
		if err != nil {
			return err
		}
		newContent, err := tempEditor(&noteData.Content)
		if err != nil {
			return err
		}
		noteData.Content = newContent
		if err := noteCmd.store.Save(noteData); err != nil {
			return err
		}
	case "delete":
		switch noteArgs[0] {
		case "tag":
			if len(noteArgs[1:]) < 2 {
				return errors.New("missing operand")
			}
			noteData, err := noteCmd.store.Load(noteArgs[1])
			if err != nil {
				return err
			}
			if err := noteData.RemoveTag(noteArgs[2]); err != nil {
				return err
			}
			if err := noteCmd.store.Save(noteData); err != nil {
				return err
			}
		case "link":
			if len(noteArgs[1:]) < 2 {
				return errors.New("missing operand")
			}
			noteData, err := noteCmd.store.Load(noteArgs[1])
			if err != nil {
				return err
			}
			if err := noteData.RemoveLink(noteArgs[2]); err != nil {
				return err
			}
			if err := noteCmd.store.Save(noteData); err != nil {
				return err
			}
		}
		if err := noteCmd.store.Delete(noteArgs[0]); err != nil {
			return err
		}
	case "help":
		noteCmd.Help()
	}
	return nil
}

func tempEditor(content *string) (string, error) {
	tempFile, err := os.CreateTemp("", "pkm-*.md")
	if err != nil {
		return "", err
	}
	defer os.Remove(tempFile.Name())

	if content != nil {
		tempFile.WriteString(*content)
	}

	tempFile.Close()

	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vi"
	}
	cmd := exec.Command(editor, tempFile.Name())
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		return "", err
	}

	newContent, err := os.ReadFile(tempFile.Name())
	if err != nil {
		return "", err
	}
	return string(newContent), nil
}
