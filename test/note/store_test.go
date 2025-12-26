package note_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/sahay-shashank/personal-knowledge-manager/internal/crypt"
	"github.com/sahay-shashank/personal-knowledge-manager/internal/note"
)

// setupTestUser initializes a test user in .crypt
func setupTestUser(t *testing.T, tmpDir, username, password string) {
	if err := crypt.InitUser(tmpDir, username, password); err != nil {
		t.Fatalf("Failed to initialize test user: %v", err)
	}
}

func TestStoreInitialize(t *testing.T) {
	tmpDir := t.TempDir()

	store := note.InitStore(tmpDir)

	if store == nil {
		t.Fatal("Store is nil")
	}

	if store.StoreLocation != tmpDir {
		t.Errorf("Store location mismatch: got %q, want %q", store.StoreLocation, tmpDir)
	}
}

func TestSaveAndLoad(t *testing.T) {
	tmpDir := t.TempDir()
	store := note.InitStore(tmpDir)

	username := "testuser"
	password := "password"

	// ADD THIS LINE
	setupTestUser(t, tmpDir, username, password)

	kp, _ := crypt.NewKeyProvider(tmpDir, username, password)

	n := note.NewNote("Test Title", "Test Content")
	linkID := "550e8400-e29b-41d4-a716-446655440000"
	n.AddLink(linkID)
	n.AddTag("test-tag")

	err := store.Save(n, username, kp)
	if err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	expectedPath := filepath.Join(tmpDir, username, n.Id+".pkm")
	if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
		t.Errorf("Note file not created: %s", expectedPath)
	}

	loaded, err := store.Load(n.Id, username, kp)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if loaded.Id != n.Id {
		t.Errorf("ID mismatch: got %q, want %q", loaded.Id, n.Id)
	}

	if loaded.Title != n.Title {
		t.Errorf("Title mismatch: got %q, want %q", loaded.Title, n.Title)
	}

	if loaded.Content != n.Content {
		t.Errorf("Content mismatch: got %q, want %q", loaded.Content, n.Content)
	}
}

func TestSaveMultipleUsers(t *testing.T) {
	tmpDir := t.TempDir()
	store := note.InitStore(tmpDir)

	users := []string{"user1", "user2", "user3"}
	password := "testpass"

	// Initialize all users first
	for _, username := range users {
		setupTestUser(t, tmpDir, username, password)
	}

	// Save notes for each user
	for _, username := range users {
		kp, err := crypt.NewKeyProvider(tmpDir, username, password)
		if err != nil {
			t.Fatalf("NewKeyProvider for %s failed: %v", username, err)
		}

		n := note.NewNote("Title for "+username, "Content for "+username)

		err = store.Save(n, username, kp) // ← FIX: Correct parameter order
		if err != nil {
			t.Fatalf("Save for %s failed: %v", username, err)
		}

		// Verify user directory exists
		userDir := filepath.Join(tmpDir, username)
		if _, err := os.Stat(userDir); os.IsNotExist(err) {
			t.Errorf("User directory not created: %s", userDir)
		}
	}

	// Verify each user has their notes
	for _, username := range users {
		_, err := crypt.NewKeyProvider(tmpDir, username, password)
		if err != nil {
			t.Fatalf("NewKeyProvider for %s failed: %v", username, err)
		}

		files, err := os.ReadDir(filepath.Join(tmpDir, username))
		if err != nil {
			t.Fatalf("ReadDir failed for %s: %v", username, err)
		}

		// Verify at least one .pkm file exists
		hasNoteFile := false
		for _, f := range files {
			if filepath.Ext(f.Name()) == ".pkm" {
				hasNoteFile = true
				break
			}
		}

		if !hasNoteFile {
			t.Errorf("No .pkm files found for %s", username)
		}
	}
}

func TestDelete(t *testing.T) {
	tmpDir := t.TempDir()
	store := note.InitStore(tmpDir)

	username := "testuser"
	password := "password"

	// ADD THIS LINE
	setupTestUser(t, tmpDir, username, password)

	kp, _ := crypt.NewKeyProvider(tmpDir, username, password)

	n := note.NewNote("To Delete", "Delete me")
	store.Save(n, username, kp)

	notePath := filepath.Join(tmpDir, username, n.Id+".pkm")
	if _, err := os.Stat(notePath); os.IsNotExist(err) {
		t.Fatal("Note file not found before delete")
	}

	err := store.Delete(n.Id, username)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	if _, err := os.Stat(notePath); !os.IsNotExist(err) {
		t.Fatal("Note file still exists after delete")
	}
}

func TestLoadNonExistent(t *testing.T) {
	tmpDir := t.TempDir()
	store := note.InitStore(tmpDir)

	username := "testuser"
	password := "password"

	// ADD THIS LINE
	setupTestUser(t, tmpDir, username, password)

	kp, _ := crypt.NewKeyProvider(tmpDir, username, password)

	_, err := store.Load("non-existent-id", username, kp)
	if err == nil {
		t.Fatal("Load should fail for non-existent note")
	}
}

func TestSaveEncryptedData(t *testing.T) {
	tmpDir := t.TempDir()
	store := note.InitStore(tmpDir)

	username := "testuser"
	password := "password"

	// ADD THIS LINE
	setupTestUser(t, tmpDir, username, password)

	kp, _ := crypt.NewKeyProvider(tmpDir, username, password)

	n := note.NewNote("Secret", "Encrypted Content")
	store.Save(n, username, kp)

	filePath := filepath.Join(tmpDir, username, n.Id+".pkm")
	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("ReadFile failed: %v", err)
	}

	if string(fileContent[:4]) != "PKM\n" {
		t.Fatal("Magic header missing")
	}

	encryptedPart := fileContent[4:]
	var testNote note.Note
	err = json.Unmarshal(encryptedPart, &testNote)
	if err == nil {
		t.Fatal("Encrypted data should not be readable as plaintext JSON")
	}
}

func TestUpdateNote(t *testing.T) {
	tmpDir := t.TempDir()
	store := note.InitStore(tmpDir)

	username := "testuser"
	password := "password"

	// ADD THIS LINE
	setupTestUser(t, tmpDir, username, password)

	kp, _ := crypt.NewKeyProvider(tmpDir, username, password)

	n := note.NewNote("Original Title", "Original Content")
	store.Save(n, username, kp)

	n.Title = "Updated Title"
	n.Content = "Updated Content"
	store.Save(n, username, kp)

	loaded, _ := store.Load(n.Id, username, kp)
	if loaded.Title != "Updated Title" {
		t.Errorf("Title not updated: got %q", loaded.Title)
	}

	if loaded.Content != "Updated Content" {
		t.Errorf("Content not updated: got %q", loaded.Content)
	}
}
