package cli_test

import (
	"testing"

	"github.com/sahay-shashank/personal-knowledge-manager/internal/cli"
	"github.com/sahay-shashank/personal-knowledge-manager/internal/note"
)

// TestLinkCommandName tests LinkCommand.Name()
func TestLinkCommandName(t *testing.T) {
	tmpDir := t.TempDir()
	testCli := setupTestEnvironment(t, tmpDir, "testuser", "password")
	cliObj := testCli.toCli()

	linkCmd := &cli.LinkCommand{Cli: cliObj}
	if linkCmd.Name() != "link" {
		t.Errorf("Expected 'link', got %q", linkCmd.Name())
	}
}

// TestLinkCommandDescription tests LinkCommand.Description()
func TestLinkCommandDescription(t *testing.T) {
	tmpDir := t.TempDir()
	testCli := setupTestEnvironment(t, tmpDir, "testuser", "password")
	cliObj := testCli.toCli()

	linkCmd := &cli.LinkCommand{Cli: cliObj}
	desc := linkCmd.Description()
	if desc == "" {
		t.Error("Description should not be empty")
	}
}

// TestLinkCommandAdd tests adding a link between two notes
func TestLinkCommandAdd(t *testing.T) {
	tmpDir := t.TempDir()
	testCli := setupTestEnvironment(t, tmpDir, "testuser", "password")
	cliObj := testCli.toCli()
	linkCmd := &cli.LinkCommand{Cli: cliObj}

	// Create two notes
	n1 := note.NewNote("Note 1", "Content 1")
	n2 := note.NewNote("Note 2", "Content 2")
	if err := testCli.Store.Save(n1, testCli.Username, testCli.KeyProvider); err != nil {
		t.Fatalf("Failed to save note 1: %v", err)
	}
	if err := testCli.Store.Save(n2, testCli.Username, testCli.KeyProvider); err != nil {
		t.Fatalf("Failed to save note 2: %v", err)
	}

	// Add link
	err := linkCmd.Run([]string{"add", n1.Id, n2.Id})
	if err != nil {
		t.Errorf("Add link failed: %v", err)
	}

	// Verify link was added
	loaded1, err := testCli.Store.Load(n1.Id, testCli.Username, testCli.KeyProvider)
	if err != nil {
		t.Fatalf("Failed to load note 1: %v", err)
	}

	found1 := false
	for _, linkID := range loaded1.Links {
		if linkID == n2.Id {
			found1 = true
			break
		}
	}
	if !found1 {
		t.Error("Link not added to note 1")
	}

	// Verify link was added
	loaded2, err := testCli.Store.Load(n2.Id, testCli.Username, testCli.KeyProvider)
	if err != nil {
		t.Fatalf("Failed to load note 2: %v", err)
	}

	found := false
	for _, linkID := range loaded2.Links {
		if linkID == n1.Id {
			found = true
			break
		}
	}
	if !found {
		t.Error("Link not backlinked to note 2")
	}
}

// TestLinkCommandRemove tests removing a link between two notes
func TestLinkCommandRemove(t *testing.T) {
	tmpDir := t.TempDir()
	testCli := setupTestEnvironment(t, tmpDir, "testuser", "password")
	cliObj := testCli.toCli()
	linkCmd := &cli.LinkCommand{Cli: cliObj}

	// Create two notes and link them
	n1 := note.NewNote("Note 1", "Content 1")
	n2 := note.NewNote("Note 2", "Content 2")
	if err := testCli.Store.Save(n1, testCli.Username, testCli.KeyProvider); err != nil {
		t.Fatalf("Failed to save note 1: %v", err)
	}
	if err := testCli.Store.Save(n2, testCli.Username, testCli.KeyProvider); err != nil {
		t.Fatalf("Failed to save note 2: %v", err)
	}

	n1.AddLink(n2.Id)
	if err := testCli.Store.Save(n1, testCli.Username, testCli.KeyProvider); err != nil {
		t.Fatalf("Failed to save link: %v", err)
	}

	// Remove link
	err := linkCmd.Run([]string{"remove", n1.Id, n2.Id})
	if err != nil {
		t.Errorf("Remove link failed: %v", err)
	}

	// Verify link was removed
	loaded, err := testCli.Store.Load(n1.Id, testCli.Username, testCli.KeyProvider)
	if err != nil {
		t.Fatalf("Failed to load note: %v", err)
	}

	// Debug: print links
	t.Logf("Loaded links after remove: %v", loaded.Links)

	for _, linkID := range loaded.Links {
		if linkID == n2.Id {
			t.Error("Link not removed from note")
		}
	}
}

// TestLinkCommandMissingArgs tests link command fails without arguments
func TestLinkCommandMissingArgs(t *testing.T) {
	tmpDir := t.TempDir()
	testCli := setupTestEnvironment(t, tmpDir, "testuser", "password")
	cliObj := testCli.toCli()
	linkCmd := &cli.LinkCommand{Cli: cliObj}

	err := linkCmd.Run([]string{})
	if err == nil {
		t.Error("Expected error for missing arguments")
	}
}
