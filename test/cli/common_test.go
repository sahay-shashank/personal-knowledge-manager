package cli_test

import (
	"testing"

	"github.com/sahay-shashank/personal-knowledge-manager/internal/cli"
	"github.com/sahay-shashank/personal-knowledge-manager/internal/crypt"
	"github.com/sahay-shashank/personal-knowledge-manager/internal/note"
)

// TestCli is a wrapper for testing that exposes private fields
type TestCli struct {
	Store       *note.Store
	Username    string
	KeyProvider *crypt.KeyProvider
}

// setupTestEnvironment creates a test CLI with initialized user
func setupTestEnvironment(t *testing.T, tmpDir, username, password string) *TestCli {
	// Initialize user
	if err := crypt.InitUser(tmpDir, username, password); err != nil {
		t.Fatalf("Failed to initialize test user: %v", err)
	}

	// Create key provider
	kp, err := crypt.NewKeyProvider(tmpDir, username, password)
	if err != nil {
		t.Fatalf("Failed to create key provider: %v", err)
	}

	// Return test CLI wrapper
	return &TestCli{
		Store:       note.InitStore(tmpDir),
		Username:    username,
		KeyProvider: kp,
	}
}

// toCli converts TestCli to actual Cli for command testing
func (tc *TestCli) toCli() *cli.Cli {
	cliObj := &cli.Cli{}
	cliObj.SetStore(tc.Store)
	cliObj.SetUsername(tc.Username)
	cliObj.SetKeyProvider(tc.KeyProvider)
	return cliObj
}

// TestCLIStoreGetter tests the Store getter
func TestCLIStoreGetter(t *testing.T) {
	tmpDir := t.TempDir()
	testCli := setupTestEnvironment(t, tmpDir, "testuser", "password")
	cliObj := testCli.toCli()

	store := cliObj.GetStore()
	if store == nil {
		t.Error("CLI Store should not be nil")
	}

	if store.StoreLocation != tmpDir {
		t.Errorf("Store location mismatch: got %q, want %q", store.StoreLocation, tmpDir)
	}
}

// TestCLIUsernameGetter tests the Username getter
func TestCLIUsernameGetter(t *testing.T) {
	tmpDir := t.TempDir()
	testCli := setupTestEnvironment(t, tmpDir, "testuser", "password")
	cliObj := testCli.toCli()

	username := cliObj.GetUsername()
	if username != "testuser" {
		t.Errorf("Username mismatch: got %q, want %q", username, "testuser")
	}
}

// TestCLIKeyProviderGetter tests the KeyProvider getter
func TestCLIKeyProviderGetter(t *testing.T) {
	tmpDir := t.TempDir()
	testCli := setupTestEnvironment(t, tmpDir, "testuser", "password")
	cliObj := testCli.toCli()

	kp := cliObj.GetKeyProvider()
	if kp == nil {
		t.Error("CLI KeyProvider should not be nil")
	}
}

// TestMultipleUsers tests that multiple users can store notes independently
func TestMultipleUsers(t *testing.T) {
	tmpDir := t.TempDir()

	// Create two users
	user1 := setupTestEnvironment(t, tmpDir, "user1", "pass1")
	user2 := setupTestEnvironment(t, tmpDir, "user2", "pass2")

	// User 1 creates a note
	n1 := note.NewNote("User 1 Note", "Content 1")
	if err := user1.Store.Save(n1, user1.Username, user1.KeyProvider); err != nil {
		t.Fatalf("Failed to save note for user1: %v", err)
	}

	// User 2 creates a note
	n2 := note.NewNote("User 2 Note", "Content 2")
	if err := user2.Store.Save(n2, user2.Username, user2.KeyProvider); err != nil {
		t.Fatalf("Failed to save note for user2: %v", err)
	}

	// User 1 should see their note but not user2's
	list1, err := user1.Store.List(user1.Username, user1.KeyProvider)
	if err != nil {
		t.Fatalf("Failed to list notes for user1: %v", err)
	}

	if len(list1) != 1 {
		t.Errorf("User 1 should have 1 note, got %d", len(list1))
	}

	// User 2 should see their note but not user1's
	list2, err := user2.Store.List(user2.Username, user2.KeyProvider)
	if err != nil {
		t.Fatalf("Failed to list notes for user2: %v", err)
	}

	if len(list2) != 1 {
		t.Errorf("User 2 should have 1 note, got %d", len(list2))
	}
}
