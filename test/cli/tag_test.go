package cli_test

import (
	"testing"

	"github.com/sahay-shashank/personal-knowledge-manager/internal/cli"
	"github.com/sahay-shashank/personal-knowledge-manager/internal/note"
)

// TestTagCommandName tests TagCommand.Name()
func TestTagCommandName(t *testing.T) {
	tmpDir := t.TempDir()
	testCli := setupTestEnvironment(t, tmpDir, "testuser", "password")
	cliObj := testCli.toCli()

	tagCmd := &cli.TagCommand{Cli: cliObj}
	if tagCmd.Name() != "tag" {
		t.Errorf("Expected 'tag', got %q", tagCmd.Name())
	}
}

// TestTagCommandDescription tests TagCommand.Description()
func TestTagCommandDescription(t *testing.T) {
	tmpDir := t.TempDir()
	testCli := setupTestEnvironment(t, tmpDir, "testuser", "password")
	cliObj := testCli.toCli()

	tagCmd := &cli.TagCommand{Cli: cliObj}
	desc := tagCmd.Description()
	if desc == "" {
		t.Error("Description should not be empty")
	}
}

// TestTagCommandAdd tests adding a tag to a note
func TestTagCommandAdd(t *testing.T) {
	tmpDir := t.TempDir()
	testCli := setupTestEnvironment(t, tmpDir, "testuser", "password")
	cliObj := testCli.toCli()
	tagCmd := &cli.TagCommand{Cli: cliObj}

	// Create a note
	n := note.NewNote("Test Note", "Test Content")
	if err := testCli.Store.Save(n, testCli.Username, testCli.KeyProvider); err != nil {
		t.Fatalf("Failed to save note: %v", err)
	}

	// Add tag
	err := tagCmd.Run([]string{"add", n.Id, "test-tag"})
	if err != nil {
		t.Errorf("Add tag failed: %v", err)
	}

	// Verify tag was added
	loaded, err := testCli.Store.Load(n.Id, testCli.Username, testCli.KeyProvider)
	if err != nil {
		t.Fatalf("Failed to load note: %v", err)
	}

	found := false
	for _, tag := range loaded.Tags {
		if tag == "test-tag" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Tag not added to note")
	}
}

// TestTagCommandDelete tests removing a tag from a note
func TestTagCommandDelete(t *testing.T) {
	tmpDir := t.TempDir()
	testCli := setupTestEnvironment(t, tmpDir, "testuser", "password")
	cliObj := testCli.toCli()
	tagCmd := &cli.TagCommand{Cli: cliObj}

	// Create a note with a tag
	n := note.NewNote("Test Note", "Test Content")
	n.AddTag("test-tag")
	if err := testCli.Store.Save(n, testCli.Username, testCli.KeyProvider); err != nil {
		t.Fatalf("Failed to save note: %v", err)
	}

	// Remove tag
	err := tagCmd.Run([]string{"remove", n.Id, "test-tag"})
	if err != nil {
		t.Errorf("Delete tag failed: %v", err)
	}

	// Verify tag was removed
	loaded, err := testCli.Store.Load(n.Id, testCli.Username, testCli.KeyProvider)
	if err != nil {
		t.Fatalf("Failed to load note: %v", err)
	}

	for _, tag := range loaded.Tags {
		if tag == "test-tag" {
			t.Error("Tag not removed from note")
		}
	}
}

// TestTagCommandAddMultipleTags tests adding multiple tags to a note
func TestTagCommandAddMultipleTags(t *testing.T) {
	tmpDir := t.TempDir()
	testCli := setupTestEnvironment(t, tmpDir, "testuser", "password")
	cliObj := testCli.toCli()
	tagCmd := &cli.TagCommand{Cli: cliObj}

	// Create a note
	n := note.NewNote("Test Note", "Test Content")
	if err := testCli.Store.Save(n, testCli.Username, testCli.KeyProvider); err != nil {
		t.Fatalf("Failed to save note: %v", err)
	}

	// Add multiple tags
	tags := []string{"tag1", "tag2", "tag3"}
	for _, tag := range tags {
		err := tagCmd.Run([]string{"add", n.Id, tag})
		if err != nil {
			t.Errorf("Add tag failed for %s: %v", tag, err)
		}
	}

	// Verify all tags were added
	loaded, err := testCli.Store.Load(n.Id, testCli.Username, testCli.KeyProvider)
	if err != nil {
		t.Fatalf("Failed to load note: %v", err)
	}

	if len(loaded.Tags) != 3 {
		t.Errorf("Expected 3 tags, got %d", len(loaded.Tags))
	}
}
