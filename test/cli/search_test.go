package cli_test

import (
	"testing"

	"github.com/sahay-shashank/personal-knowledge-manager/internal/cli"
	"github.com/sahay-shashank/personal-knowledge-manager/internal/note"
)

// TestSearchCommandName tests SearchCommand.Name()
func TestSearchCommandName(t *testing.T) {
	tmpDir := t.TempDir()
	testCli := setupTestEnvironment(t, tmpDir, "testuser", "password")
	cliObj := testCli.toCli()

	searchCmd := &cli.SearchCommand{Cli: cliObj}
	if searchCmd.Name() != "search" {
		t.Errorf("Expected 'search', got %q", searchCmd.Name())
	}
}

// TestSearchCommandDescription tests SearchCommand.Description()
func TestSearchCommandDescription(t *testing.T) {
	tmpDir := t.TempDir()
	testCli := setupTestEnvironment(t, tmpDir, "testuser", "password")
	cliObj := testCli.toCli()

	searchCmd := &cli.SearchCommand{Cli: cliObj}
	desc := searchCmd.Description()
	if desc == "" {
		t.Error("Description should not be empty")
	}
}

// TestSearchCommandByTag tests searching notes by tag
func TestSearchCommandByTag(t *testing.T) {
	tmpDir := t.TempDir()
	testCli := setupTestEnvironment(t, tmpDir, "testuser", "password")
	cliObj := testCli.toCli()
	searchCmd := &cli.SearchCommand{Cli: cliObj}

	// Create notes with tags
	n1 := note.NewNote("Go Note", "Go content")
	n1.AddTag("go")
	if err := testCli.Store.Save(n1, testCli.Username, testCli.KeyProvider); err != nil {
		t.Fatalf("Failed to save note: %v", err)
	}

	n2 := note.NewNote("Rust Note", "Rust content")
	n2.AddTag("rust")
	if err := testCli.Store.Save(n2, testCli.Username, testCli.KeyProvider); err != nil {
		t.Fatalf("Failed to save note: %v", err)
	}

	// Search by tag
	err := searchCmd.Run([]string{"tag", "go"})
	if err != nil {
		t.Errorf("Search by tag failed: %v", err)
	}
}

// TestSearchCommandByKeyword tests searching notes by keyword
func TestSearchCommandByKeyword(t *testing.T) {
	tmpDir := t.TempDir()
	testCli := setupTestEnvironment(t, tmpDir, "testuser", "password")
	cliObj := testCli.toCli()
	searchCmd := &cli.SearchCommand{Cli: cliObj}

	// Create notes with keywords
	n1 := note.NewNote("Go Programming", "Learn Go language")
	if err := testCli.Store.Save(n1, testCli.Username, testCli.KeyProvider); err != nil {
		t.Fatalf("Failed to save note: %v", err)
	}

	// Search by keyword
	err := searchCmd.Run([]string{"keyword", "Programming"})
	if err != nil {
		t.Errorf("Search by keyword failed: %v", err)
	}
}

// TestSearchCommandMissingArgs tests search command fails without arguments
func TestSearchCommandMissingArgs(t *testing.T) {
	tmpDir := t.TempDir()
	testCli := setupTestEnvironment(t, tmpDir, "testuser", "password")
	cliObj := testCli.toCli()
	searchCmd := &cli.SearchCommand{Cli: cliObj}

	err := searchCmd.Run([]string{})
	if err == nil {
		t.Error("Expected error for missing arguments")
	}
}

// TestSearchCommandMissingSearchTerm tests search command fails without search term
func TestSearchCommandMissingSearchTerm(t *testing.T) {
	tmpDir := t.TempDir()
	testCli := setupTestEnvironment(t, tmpDir, "testuser", "password")
	cliObj := testCli.toCli()
	searchCmd := &cli.SearchCommand{Cli: cliObj}

	err := searchCmd.Run([]string{"tag"})
	if err == nil {
		t.Error("Expected error when search type has no terms")
	}
}

// TestSearchCommandMultipleTags tests searching with multiple tags
func TestSearchCommandMultipleTags(t *testing.T) {
	tmpDir := t.TempDir()
	testCli := setupTestEnvironment(t, tmpDir, "testuser", "password")
	cliObj := testCli.toCli()
	searchCmd := &cli.SearchCommand{Cli: cliObj}

	// Create a note with multiple tags
	n := note.NewNote("Multi-tag Note", "Content with multiple tags")
	n.AddTag("go")
	n.AddTag("programming")
	if err := testCli.Store.Save(n, testCli.Username, testCli.KeyProvider); err != nil {
		t.Fatalf("Failed to save note: %v", err)
	}

	// Search by multiple tags
	err := searchCmd.Run([]string{"tag", "go", "programming"})
	if err != nil {
		t.Errorf("Search by multiple tags failed: %v", err)
	}
}
