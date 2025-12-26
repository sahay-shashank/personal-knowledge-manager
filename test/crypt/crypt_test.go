package crypto_test

import (
	"testing"

	"github.com/sahay-shashank/personal-knowledge-manager/internal/crypt"
)

// setupTestUser initializes a test user in .crypt
func setupTestUser(t *testing.T, tmpDir, username, password string) {
	if err := crypt.InitUser(tmpDir, username, password); err != nil {
		t.Fatalf("Failed to initialize test user: %v", err)
	}
}

func TestNewKeyProvider(t *testing.T) {
	tmpDir := t.TempDir()
	username := "testuser"
	password := "password123"

	// Initialize user first
	setupTestUser(t, tmpDir, username, password)

	// Now create KeyProvider
	kp, err := crypt.NewKeyProvider(tmpDir, username, password)
	if err != nil {
		t.Fatalf("NewKeyProvider failed: %v", err)
	}

	if kp == nil {
		t.Fatal("KeyProvider is nil")
	}

	if len(kp.DEK()) != crypt.DEKSize {
		t.Errorf("DEK size mismatch: got %d, want %d", len(kp.DEK()), crypt.DEKSize)
	}
}

func TestEncrypt(t *testing.T) {
	tmpDir := t.TempDir()
	username := "testuser"
	password := "password"

	setupTestUser(t, tmpDir, username, password)
	kp, _ := crypt.NewKeyProvider(tmpDir, username, password)

	plaintext := []byte("Hello, World!")

	encrypted, err := kp.Encrypt(plaintext)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	if len(encrypted) <= len(plaintext) {
		t.Errorf("Encrypted size too small: got %d, want > %d", len(encrypted), len(plaintext))
	}

	if string(encrypted) == string(plaintext) {
		t.Fatal("Encrypted equals plaintext (not encrypted)")
	}
}

func TestDecrypt(t *testing.T) {
	tmpDir := t.TempDir()
	username := "testuser"
	password := "password"

	setupTestUser(t, tmpDir, username, password)
	kp, _ := crypt.NewKeyProvider(tmpDir, username, password)

	plaintext := []byte("Hello, World!")

	encrypted, err := kp.Encrypt(plaintext)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	decrypted, err := kp.Decrypt(encrypted)
	if err != nil {
		t.Fatalf("Decrypt failed: %v", err)
	}

	if string(decrypted) != string(plaintext) {
		t.Errorf("Decrypted mismatch: got %q, want %q", decrypted, plaintext)
	}
}

func TestEncryptDecryptRoundtrip(t *testing.T) {
	tmpDir := t.TempDir()
	username := "user1"
	password := "pass"

	setupTestUser(t, tmpDir, username, password)
	kp, _ := crypt.NewKeyProvider(tmpDir, username, password)

	tests := []struct {
		name      string
		plaintext string
	}{
		{"empty", ""},
		{"short", "x"},
		{"normal", "This is a test note with some content."},
		{"multiline", "Line 1\nLine 2\nLine 3"},
		{"special chars", "!@#$%^&*()_+-=[]{}|;:',.<>?/"},
		{"json", `{"title":"test","content":"data"}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			plaintext := []byte(tt.plaintext)
			encrypted, err := kp.Encrypt(plaintext)
			if err != nil {
				t.Fatalf("Encrypt failed: %v", err)
			}

			decrypted, err := kp.Decrypt(encrypted)
			if err != nil {
				t.Fatalf("Decrypt failed: %v", err)
			}

			if string(decrypted) != tt.plaintext {
				t.Errorf("Roundtrip failed: got %q, want %q", decrypted, tt.plaintext)
			}
		})
	}
}

func TestDecryptWithWrongPassword(t *testing.T) {
	tmpDir := t.TempDir()
	username := "testuser"

	setupTestUser(t, tmpDir, username, "correctpassword")

	_, err := crypt.NewKeyProvider(tmpDir, username, "wrongpassword")
	if err == nil {
		t.Fatal("Decrypt with wrong password should fail")
	}
}

func TestDecryptWithWrongDEK(t *testing.T) {
	tmpDir := t.TempDir()

	setupTestUser(t, tmpDir, "user1", "pass1")
	setupTestUser(t, tmpDir, "user2", "pass2")

	kp1, _ := crypt.NewKeyProvider(tmpDir, "user1", "pass1")
	kp2, _ := crypt.NewKeyProvider(tmpDir, "user2", "pass2")

	plaintext := []byte("Secret message")
	encrypted, _ := kp1.Encrypt(plaintext)

	_, err := kp2.Decrypt(encrypted)
	if err == nil {
		t.Fatal("Decrypt with wrong DEK should fail")
	}
}

func TestDecryptCorruptedData(t *testing.T) {
	tmpDir := t.TempDir()
	setupTestUser(t, tmpDir, "testuser", "password")
	kp, _ := crypt.NewKeyProvider(tmpDir, "testuser", "password")

	tests := []struct {
		name string
		data []byte
	}{
		{"empty", []byte{}},
		{"too short", []byte{1, 2, 3}},
		{"invalid", []byte("not encrypted data")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := kp.Decrypt(tt.data)
			if err == nil {
				t.Errorf("Decrypt should fail for %s", tt.name)
			}
		})
	}
}

func TestInitUser(t *testing.T) {
	tmpDir := t.TempDir()
	username := "newuser"
	password := "password123"

	err := crypt.InitUser(tmpDir, username, password)
	if err != nil {
		t.Fatalf("InitUser failed: %v", err)
	}

	// Verify user can be loaded
	kp, err := crypt.NewKeyProvider(tmpDir, username, password)
	if err != nil {
		t.Fatalf("Failed to load initialized user: %v", err)
	}

	if kp == nil {
		t.Fatal("KeyProvider is nil after init")
	}
}
