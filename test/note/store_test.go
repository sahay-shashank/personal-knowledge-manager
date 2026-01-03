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

func TestListSorting(t *testing.T) {
	tmpDir := t.TempDir()
	store := note.InitStore(tmpDir)
	username := "listtest"
	password := "pass"
	setupTestUser(t, tmpDir, username, password)
	kp, _ := crypt.NewKeyProvider(tmpDir, username, password)

	// Save unsorted notes
	store.Save(&note.Note{Id: "z", Title: "Zebra"}, username, kp)
	store.Save(&note.Note{Id: "a", Title: "Apple"}, username, kp)
	store.Save(&note.Note{Id: "m", Title: "Monkey"}, username, kp)

	summaries, err := store.List(username, kp)
	if err != nil {
		t.Fatal(err)
	}
	if len(summaries) != 3 {
		t.Fatalf("want 3, got %d", len(summaries))
	}
	titles := []string{summaries[0].Title, summaries[1].Title, summaries[2].Title}
	if titles[0] != "Apple" || titles[1] != "Monkey" || titles[2] != "Zebra" {
		t.Errorf("want [Apple Monkey Zebra], got %v", titles)
	}
}

func TestListSkipsNonPKM(t *testing.T) {
	tmpDir := t.TempDir()
	store := note.InitStore(tmpDir)
	username := "skiptest"
	password := "pass"
	setupTestUser(t, tmpDir, username, password)
	kp, _ := crypt.NewKeyProvider(tmpDir, username, password)

	store.Save(&note.Note{Id: "valid", Title: "Valid"}, username, kp)

	// Add junk files
	junkPath := filepath.Join(tmpDir, username, "junk.txt")
	os.WriteFile(junkPath, []byte("junk"), 0644)
	dirPath := filepath.Join(tmpDir, username, "subdir")
	os.Mkdir(dirPath, 0755)

	summaries, err := store.List(username, kp)
	if err != nil {
		t.Fatal(err)
	}
	if len(summaries) != 1 {
		t.Errorf("want 1 .pkm file only, got %d", len(summaries))
	}
}

func TestListCorruptedHeader(t *testing.T) {
	tmpDir := t.TempDir()
	store := note.InitStore(tmpDir)
	username := "corrupt"
	password := "pass"
	setupTestUser(t, tmpDir, username, password)
	kp, _ := crypt.NewKeyProvider(tmpDir, username, password)

	// Valid note
	store.Save(&note.Note{Id: "good", Title: "Good"}, username, kp)

	// Corrupted: wrong header
	corruptPath := filepath.Join(tmpDir, username, "bad.pkm")
	os.WriteFile(corruptPath, []byte("BAD\n..."), 0644)

	summaries, err := store.List(username, kp)
	if err != nil {
		t.Errorf("want no error (corrupted file should be skipped), got %v", err)
	}
	// Still lists good notes
	if len(summaries) != 1 || summaries[0].Title != "Good" {
		t.Errorf("want 1 good summary, got %v", summaries)
	}
}

func TestSearchNoIndex(t *testing.T) {
	tmpDir := t.TempDir()
	store := note.InitStore(tmpDir)
	username := "searchtest"
	password := "pass"
	setupTestUser(t, tmpDir, username, password)
	kp, _ := crypt.NewKeyProvider(tmpDir, username, password)

	// Save some notes (creates index)
	store.Save(&note.Note{Id: "n1", Title: "Go lang", Tags: []string{"go"}}, username, kp)

	matches, err := store.Search("tag", []string{"rust"}, username, kp)
	if err != nil {
		t.Fatal(err)
	}
	if len(matches) != 0 {
		t.Errorf("want no matches for missing tag, got %v", matches)
	}
}

func TestIndexUpdates(t *testing.T) {
	tmpDir := t.TempDir()
	store := note.InitStore(tmpDir)
	username := "indextest"
	password := "pass"
	setupTestUser(t, tmpDir, username, password)
	kp, _ := crypt.NewKeyProvider(tmpDir, username, password)

	note1 := &note.Note{Id: "one", Title: "hello world", Tags: []string{"test"}}
	store.Save(note1, username, kp)

	// Check index created
	indexPath := filepath.Join(tmpDir, username, ".index.pkm")
	indexData, _ := os.ReadFile(indexPath)

	// Decrypt the index
	decryptedData, err := kp.Decrypt(indexData[4:]) // skip PKM header
	if err != nil {
		t.Fatalf("Failed to decrypt index: %v", err)
	}

	var idx note.Index
	json.Unmarshal(decryptedData, &idx)

	if len(idx.TagIndex["test"]) != 1 || idx.TagIndex["test"][0] != "one" {
		t.Error("tag index not updated")
	}
	if len(idx.KeywordIndex["world"]) != 1 {
		t.Error("keyword index not updated")
	}
}
