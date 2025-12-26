package note_test

import (
	"testing"

	"github.com/sahay-shashank/personal-knowledge-manager/internal/note"
)

func TestNewNote(t *testing.T) {
	title := "Test Note"
	content := "Test content"

	n := note.NewNote(title, content)

	if n == nil {
		t.Fatal("Note is nil")
	}

	if n.Title != title {
		t.Errorf("Title mismatch: got %q, want %q", n.Title, title)
	}

	if n.Content != content {
		t.Errorf("Content mismatch: got %q, want %q", n.Content, content)
	}

	if n.Id == "" {
		t.Fatal("ID is empty")
	}

	if len(n.Links) != 0 {
		t.Errorf("Links should be empty: got %v", n.Links)
	}

	if len(n.Tags) != 0 {
		t.Errorf("Tags should be empty: got %v", n.Tags)
	}
}

func TestAddLink(t *testing.T) {
	n := note.NewNote("Test", "Content")
	targetID := "target-123"

	err := n.AddLink(targetID)
	if err != nil {
		t.Fatalf("AddLink failed: %v", err)
	}

	if len(n.Links) != 1 || n.Links[0] != targetID {
		t.Errorf("Link not added correctly: got %v", n.Links)
	}

	// Adding same link again should error
	err = n.AddLink(targetID)
	if err == nil {
		t.Fatal("Adding duplicate link should fail")
	}
}

func TestRemoveLink(t *testing.T) {
	n := note.NewNote("Test", "Content")
	link1 := "link-1"
	link2 := "link-2"

	n.AddLink(link1)
	n.AddLink(link2)

	err := n.RemoveLink(link1)
	if err != nil {
		t.Fatalf("RemoveLink failed: %v", err)
	}

	if len(n.Links) != 1 || n.Links[0] != link2 {
		t.Errorf("Link not removed correctly: got %v", n.Links)
	}

	err = n.RemoveLink("non-existent")
	if err == nil {
		t.Fatal("Removing non-existent link should fail")
	}
}

func TestAddTag(t *testing.T) {
	n := note.NewNote("Test", "Content")

	err := n.AddTag("tag1, tag2, tag3")
	if err != nil {
		t.Fatalf("AddTag failed: %v", err)
	}

	if len(n.Tags) != 3 {
		t.Errorf("Tags count mismatch: got %d, want 3", len(n.Tags))
	}

	// Case insensitivity
	err = n.AddTag("TAG1")
	if err == nil {
		t.Fatal("Adding duplicate tag (case insensitive) should fail")
	}
}

func TestRemoveTag(t *testing.T) {
	n := note.NewNote("Test", "Content")
	n.AddTag("tag1, tag2, tag3")

	err := n.RemoveTag("tag1")
	if err != nil {
		t.Fatalf("RemoveTag failed: %v", err)
	}

	if len(n.Tags) != 2 {
		t.Errorf("Tags count mismatch: got %d, want 2", len(n.Tags))
	}
}
