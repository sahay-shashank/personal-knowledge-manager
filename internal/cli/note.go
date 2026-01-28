package cli

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/sahay-shashank/personal-knowledge-manager/internal/note"
)

type NoteCommand struct {
	*Cli
}

func (noteCmd *NoteCommand) Name() string {
	return "note"
}

func (noteCmd *NoteCommand) Description() string {
	return "Create, read, edit, and delete encrypted notes"
}

func (noteCmd *NoteCommand) Run(args []string) error {
	if len(args) < 1 {
		noteCmd.Help()
		return errors.New("missing arguments")
	}
	cmd := args[0]
	noteArgs := args[1:]
	switch cmd {
	case "new":
		if len(noteArgs) < 1 {
			return errors.New("usage: note new <title>")
		}

		content, err := tempEditor(nil)
		if err != nil {
			return err
		}
		if content == "" {
			return errors.New("no content")
		}
		noteData := note.NewNote(strings.Join(noteArgs, " "), content)
		return noteCmd.store.Save(noteData, noteCmd.username, noteCmd.keyProvider)

	case "edit":
		if len(noteArgs) < 1 {
			return errors.New("usage: note new <title>")
		}

		noteData, err := noteCmd.store.Load(noteArgs[0], noteCmd.username, noteCmd.keyProvider)
		if err != nil {
			return err
		}
		newContent, err := tempEditor(&noteData.Content)
		if err != nil {
			return err
		}
		noteData.Content = newContent
		return noteCmd.store.Save(noteData, noteCmd.username, noteCmd.keyProvider)

	case "delete":
		if len(noteArgs) < 1 {
			return errors.New("usage: note delete <id>")
		}
		return noteCmd.store.Delete(noteArgs[0], noteCmd.username)

	case "list":
		return noteCmd.printList()

	case "help":
		noteCmd.Help()

	default:
		return fmt.Errorf("unknown subcommand: %s", cmd)
	}
	return nil
}

func (noteCmd *NoteCommand) printList() error {
	noteSummaryList, err := noteCmd.store.List(noteCmd.username, noteCmd.keyProvider)
	if len(noteSummaryList) == 0 {
		fmt.Println("No Notes found!")
		return nil
	}
	if err != nil {
		return err
	}
	maxUID, maxTitle, maxTags := 3, 5, 4
	for _, s := range noteSummaryList {
		maxUID = max(maxUID, len(s.Id))
		maxTitle = max(maxTitle, len(s.Title))
		maxTags = max(maxTags, len(strings.Join(s.Tags, ",")))
	}

	dashUID := strings.Repeat("-", maxUID)
	dashTitle := strings.Repeat("-", maxTitle)
	dashTags := strings.Repeat("-", maxTags)
	separator := fmt.Sprintf("%s\t%s\t%s", dashUID, dashTitle, dashTags)

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "UID\tTITLE\tTAGS")
	fmt.Fprintln(w, separator)
	for _, noteSummary := range noteSummaryList {
		tags := strings.Join(noteSummary.Tags, ",")
		fmt.Fprintf(w, "%s\t%s\t%s\n", noteSummary.Id, noteSummary.Title, tags)
	}

	return w.Flush()
}
