package cli_test

import (
	"testing"

	"github.com/sahay-shashank/personal-knowledge-manager/internal/cli"
	"github.com/sahay-shashank/personal-knowledge-manager/internal/note"
)

// TestNoteCommandName tests NoteCommand.Name()
func TestNoteCommandName(t *testing.T) {
	tmpDir := t.TempDir()
	testCli := setupTestEnvironment(t, tmpDir, "testuser", "password")
	cliObj := testCli.toCli()

	noteCmd := &cli.NoteCommand{Cli: cliObj}
	if noteCmd.Name() != "note" {
		t.Errorf("Expected 'note', got %q", noteCmd.Name())
	}
}

// TestNoteCommandDescription tests NoteCommand.Description()
func TestNoteCommandDescription(t *testing.T) {
	tmpDir := t.TempDir()
	testCli := setupTestEnvironment(t, tmpDir, "testuser", "password")
	cliObj := testCli.toCli()

	noteCmd := &cli.NoteCommand{Cli: cliObj}
	desc := noteCmd.Description()
	if desc == "" {
		t.Error("Description should not be empty")
	}
	if len(desc) < 5 {
		t.Errorf("Description too short: %q", desc)
	}
}

// TestNoteCommandMissingArgs tests that note command fails without arguments
func TestNoteCommandMissingArgs(t *testing.T) {
	tmpDir := t.TempDir()
	testCli := setupTestEnvironment(t, tmpDir, "testuser", "password")
	cliObj := testCli.toCli()

	noteCmd := &cli.NoteCommand{Cli: cliObj}
	err := noteCmd.Run([]string{})
	if err == nil {
		t.Error("Expected error for missing arguments")
	}
}

// TestNoteCommandDelete tests deleting a note
func TestNoteCommandDelete(t *testing.T) {
	tmpDir := t.TempDir()
	testCli := setupTestEnvironment(t, tmpDir, "testuser", "password")
	cliObj := testCli.toCli()
	noteCmd := &cli.NoteCommand{Cli: cliObj}

	// Create a note first
	n := note.NewNote("Test Note", "Test Content")
	if err := testCli.Store.Save(n, testCli.Username, testCli.KeyProvider); err != nil {
		t.Fatalf("Failed to save note: %v", err)
	}

	// Delete the note
	err := noteCmd.Run([]string{"delete", n.Id})
	if err != nil {
		t.Errorf("Delete failed: %v", err)
	}

	// Verify it's deleted
	_, err = testCli.Store.Load(n.Id, testCli.Username, testCli.KeyProvider)
	if err == nil {
		t.Error("Note should be deleted but still exists")
	}
}

// TestNoteCommandDeleteNonExistent tests deleting a non-existent note
func TestNoteCommandDeleteNonExistent(t *testing.T) {
	tmpDir := t.TempDir()
	testCli := setupTestEnvironment(t, tmpDir, "testuser", "password")
	cliObj := testCli.toCli()
	noteCmd := &cli.NoteCommand{Cli: cliObj}

	err := noteCmd.Run([]string{"delete", "non-existent-id"})
	if err == nil {
		t.Error("Expected error when deleting non-existent note")
	}
}

// TestNoteCommandDeleteMissingId tests delete with missing note ID
func TestNoteCommandDeleteMissingId(t *testing.T) {
	tmpDir := t.TempDir()
	testCli := setupTestEnvironment(t, tmpDir, "testuser", "password")
	cliObj := testCli.toCli()
	noteCmd := &cli.NoteCommand{Cli: cliObj}

	err := noteCmd.Run([]string{"delete"})
	if err == nil {
		t.Error("Expected error when delete has no arguments")
	}
}
