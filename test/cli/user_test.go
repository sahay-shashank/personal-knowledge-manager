package cli_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/sahay-shashank/personal-knowledge-manager/internal/cli"
	"github.com/sahay-shashank/personal-knowledge-manager/internal/crypt"
	"github.com/sahay-shashank/personal-knowledge-manager/internal/note"
)

// TestUserCommandName tests UserCommand.Name()
func TestUserCommandName(t *testing.T) {
	userCmd := &cli.UserCommand{CLI: &cli.Cli{}}

	if userCmd.Name() != "user" {
		t.Errorf("Expected 'user', got %q", userCmd.Name())
	}
}

// TestUserCommandDescription tests UserCommand.Description()
func TestUserCommandDescription(t *testing.T) {
	userCmd := &cli.UserCommand{CLI: &cli.Cli{}}

	desc := userCmd.Description()
	if desc == "" {
		t.Error("Description should not be empty")
	}
}

// TestUserCommandMissingSubcommand tests user command fails without subcommand
func TestUserCommandMissingSubcommand(t *testing.T) {
	userCmd := &cli.UserCommand{CLI: &cli.Cli{}}

	err := userCmd.Run([]string{})
	if err == nil {
		t.Error("Expected error for missing subcommand")
	}
}

// TestUserCommandInitUser tests initializing a new user
func TestUserCommandInitUser(t *testing.T) {
	tmpDir := t.TempDir()
	userCmd := &cli.UserCommand{CLI: &cli.Cli{}}
	cliObj := userCmd.CLI
	cliObj.SetStore(note.InitStore(tmpDir))

	// Test init user - note: this will prompt for password in real usage
	// For testing, we use the lower-level InitUser function
	err := crypt.InitUser(tmpDir, "newuser", "password")
	if err != nil {
		t.Errorf("InitUser failed: %v", err)
	}

	// Verify user was created
	userDir := filepath.Join(tmpDir, ".crypt", "newuser")
	if _, err := os.Stat(userDir); os.IsNotExist(err) {
		t.Error("User directory not created")
	}
}

// TestUserCommandMultipleUsers tests initializing multiple users
func TestUserCommandMultipleUsers(t *testing.T) {
	tmpDir := t.TempDir()

	usernames := []string{"alice", "bob", "charlie"}
	password := "testpass"

	// Initialize multiple users
	for _, username := range usernames {
		err := crypt.InitUser(tmpDir, username, password)
		if err != nil {
			t.Errorf("Failed to init user %s: %v", username, err)
		}
	}

	// Verify all users were created
	for _, username := range usernames {
		userDir := filepath.Join(tmpDir, ".crypt", username)
		if _, err := os.Stat(userDir); os.IsNotExist(err) {
			t.Errorf("User directory not created for %s", username)
		}
	}
}

// TestUserCommandInitDuplicateUser tests initializing a user that already exists
func TestUserCommandInitDuplicateUser(t *testing.T) {
	tmpDir := t.TempDir()

	// Create first user
	err := crypt.InitUser(tmpDir, "duplicate", "password")
	if err != nil {
		t.Fatalf("Failed to create first user: %v", err)
	}

	// Try to create same user again
	err = crypt.InitUser(tmpDir, "duplicate", "password")
	if err == nil {
		t.Error("Expected error when creating duplicate user")
	}
}

// TestUserCommandInitWithPasswordAuthentication tests that user can authenticate with correct password
func TestUserCommandInitWithPasswordAuthentication(t *testing.T) {
	tmpDir := t.TempDir()
	username := "testuser"
	password := "correct-password"

	// Initialize user
	err := crypt.InitUser(tmpDir, username, password)
	if err != nil {
		t.Fatalf("Failed to init user: %v", err)
	}

	// Try to create key provider with correct password
	kp, err := crypt.NewKeyProvider(tmpDir, username, password)
	if err != nil {
		t.Errorf("Failed to authenticate with correct password: %v", err)
	}

	if kp == nil {
		t.Error("KeyProvider should not be nil")
	}
}

// TestUserCommandInitWithWrongPassword tests that authentication fails with wrong password
func TestUserCommandInitWithWrongPassword(t *testing.T) {
	tmpDir := t.TempDir()
	username := "testuser"
	password := "correct-password"

	// Initialize user
	err := crypt.InitUser(tmpDir, username, password)
	if err != nil {
		t.Fatalf("Failed to init user: %v", err)
	}

	// Try to create key provider with wrong password
	_, err = crypt.NewKeyProvider(tmpDir, username, "wrong-password")
	if err == nil {
		t.Error("Expected error when authenticating with wrong password")
	}
}
