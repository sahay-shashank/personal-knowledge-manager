package crypto_test

import (
	"encoding/json"
	"os"
	"path/filepath"
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

func TestNewKeyProviderEmptyUsername(t *testing.T) {
	tmpDir := t.TempDir()
	_, err := crypt.NewKeyProvider(tmpDir, "", "password")
	if err == nil {
		t.Fatal("NewKeyProvider with empty username should fail")
	}
}

func TestNewKeyProviderUserNotFound(t *testing.T) {
	tmpDir := t.TempDir()
	setupTestUser(t, tmpDir, "existinguser", "password")

	_, err := crypt.NewKeyProvider(tmpDir, "nonexistentuser", "password")
	if err == nil {
		t.Fatal("NewKeyProvider with non-existent user should fail")
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

func TestInitUserEmptyUsername(t *testing.T) {
	tmpDir := t.TempDir()
	err := crypt.InitUser(tmpDir, "", "password")
	if err == nil {
		t.Fatal("InitUser with empty username should fail")
	}
}

func TestInitUserEmptyPassword(t *testing.T) {
	tmpDir := t.TempDir()
	err := crypt.InitUser(tmpDir, "testuser", "")
	if err == nil {
		t.Fatal("InitUser with empty password should fail")
	}
}

func TestInitUserDuplicate(t *testing.T) {
	tmpDir := t.TempDir()
	username := "dupuser"
	password := "password"

	// Create first user
	err := crypt.InitUser(tmpDir, username, password)
	if err != nil {
		t.Fatalf("InitUser failed: %v", err)
	}

	// Try to create same user again
	err = crypt.InitUser(tmpDir, username, password)
	if err == nil {
		t.Fatal("InitUser with duplicate username should fail")
	}
}

func TestChangePassword(t *testing.T) {
	tmpDir := t.TempDir()
	username := "testuser"
	oldPassword := "oldpass"
	newPassword := "newpass"

	setupTestUser(t, tmpDir, username, oldPassword)

	// Change password
	err := crypt.ChangePassword(tmpDir, username, oldPassword, newPassword)
	if err != nil {
		t.Fatalf("ChangePassword failed: %v", err)
	}

	// Old password should fail
	_, err = crypt.NewKeyProvider(tmpDir, username, oldPassword)
	if err == nil {
		t.Fatal("Old password should not work")
	}

	// New password should work
	kp, err := crypt.NewKeyProvider(tmpDir, username, newPassword)
	if err != nil {
		t.Fatalf("NewKeyProvider with new password failed: %v", err)
	}

	if kp == nil {
		t.Fatal("KeyProvider is nil after password change")
	}
}

func TestChangePasswordWrongOldPassword(t *testing.T) {
	tmpDir := t.TempDir()
	username := "testuser"

	setupTestUser(t, tmpDir, username, "correctpass")

	err := crypt.ChangePassword(tmpDir, username, "wrongpass", "newpass")
	if err == nil {
		t.Fatal("ChangePassword with wrong old password should fail")
	}
}

func TestChangePasswordEmptyNewPassword(t *testing.T) {
	tmpDir := t.TempDir()
	username := "testuser"
	password := "password"

	setupTestUser(t, tmpDir, username, password)

	err := crypt.ChangePassword(tmpDir, username, password, "")
	if err == nil {
		t.Fatal("ChangePassword with empty new password should fail")
	}
}

func TestChangePasswordPreserveDEK(t *testing.T) {
	tmpDir := t.TempDir()
	username := "testuser"
	oldPass := "oldpass"
	newPass := "newpass"

	setupTestUser(t, tmpDir, username, oldPass)
	kp1, _ := crypt.NewKeyProvider(tmpDir, username, oldPass)
	dek1 := kp1.DEK()

	// Encrypt data with old password
	plaintext := []byte("Secret data")
	encrypted, _ := kp1.Encrypt(plaintext)

	// Change password
	crypt.ChangePassword(tmpDir, username, oldPass, newPass)

	// Load with new password and decrypt
	kp2, _ := crypt.NewKeyProvider(tmpDir, username, newPass)
	dek2 := kp2.DEK()

	// DEK should be same (preserved)
	if string(dek1) != string(dek2) {
		t.Fatal("DEK should be preserved after password change")
	}

	// Data encrypted with old key should decrypt with new key
	decrypted, err := kp2.Decrypt(encrypted)
	if err != nil {
		t.Fatalf("Decrypt after password change failed: %v", err)
	}

	if string(decrypted) != string(plaintext) {
		t.Errorf("Decrypted data mismatch: got %q, want %q", decrypted, plaintext)
	}
}

func TestGetDEK(t *testing.T) {
	tmpDir := t.TempDir()
	setupTestUser(t, tmpDir, "testuser", "password")
	kp, _ := crypt.NewKeyProvider(tmpDir, "testuser", "password")

	dek := kp.GetDEK()
	if len(dek) != crypt.DEKSize {
		t.Errorf("GetDEK size mismatch: got %d, want %d", len(dek), crypt.DEKSize)
	}

	// DEK() and GetDEK() should return same value
	if string(dek) != string(kp.DEK()) {
		t.Fatal("DEK() and GetDEK() should return same value")
	}
}

func TestExportUser(t *testing.T) {
	tmpDir := t.TempDir()
	username := "exportuser"
	password := "password123"

	setupTestUser(t, tmpDir, username, password)

	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := crypt.ExportUser(tmpDir, username, password)

	w.Close()
	os.Stdout = oldStdout

	if err != nil {
		t.Fatalf("ExportUser failed: %v", err)
	}

	// Read output
	var output []byte
	n := 1
	for n > 0 {
		b := make([]byte, 1024)
		n, _ = r.Read(b)
		output = append(output, b[:n]...)
	}

	// Verify output is valid JSON
	var entry crypt.CryptEntry
	err = json.Unmarshal(output, &entry)
	if err != nil {
		t.Fatalf("Export output is not valid JSON: %v", err)
	}

	if entry.Username != username {
		t.Errorf("Exported username mismatch: got %q, want %q", entry.Username, username)
	}
}

func TestExportUserWrongPassword(t *testing.T) {
	tmpDir := t.TempDir()
	setupTestUser(t, tmpDir, "testuser", "correctpass")

	err := crypt.ExportUser(tmpDir, "testuser", "wrongpass")
	if err == nil {
		t.Fatal("ExportUser with wrong password should fail")
	}
}

func TestExportUserNotFound(t *testing.T) {
	tmpDir := t.TempDir()
	err := crypt.ExportUser(tmpDir, "nonexistent", "password")
	if err == nil {
		t.Fatal("ExportUser with non-existent user should fail")
	}
}

func TestImportUser(t *testing.T) {
	tmpDir1 := t.TempDir()
	tmpDir2 := t.TempDir()
	username := "importuser"
	password := "password123"

	// Create user in tmpDir1
	setupTestUser(t, tmpDir1, username, password)

	// Export user data
	cryptPath1 := filepath.Join(tmpDir1, ".crypt")
	cf1, _ := crypt.ReadCryptFile(cryptPath1)
	entry1 := cf1.FindEntry(username)

	// Marshal entry
	data, _ := json.Marshal(entry1)

	// Import to tmpDir2
	err := crypt.ImportUser(tmpDir2, string(data))
	if err != nil {
		t.Fatalf("ImportUser failed: %v", err)
	}

	// Verify imported user can be loaded with same password
	kp, err := crypt.NewKeyProvider(tmpDir2, username, password)
	if err != nil {
		t.Fatalf("Failed to load imported user: %v", err)
	}

	if kp == nil {
		t.Fatal("KeyProvider is nil after import")
	}
}

func TestImportUserDuplicate(t *testing.T) {
	tmpDir := t.TempDir()
	username := "dupuser"

	setupTestUser(t, tmpDir, username, "password")

	// Try to import same user
	cryptPath := filepath.Join(tmpDir, ".crypt")
	cf, _ := crypt.ReadCryptFile(cryptPath)
	entry := cf.FindEntry(username)
	data, _ := json.Marshal(entry)

	err := crypt.ImportUser(tmpDir, string(data))
	if err == nil {
		t.Fatal("ImportUser with duplicate user should fail")
	}
}

func TestImportUserInvalidJSON(t *testing.T) {
	tmpDir := t.TempDir()
	err := crypt.ImportUser(tmpDir, "invalid json")
	if err == nil {
		t.Fatal("ImportUser with invalid JSON should fail")
	}
}

func TestImportUserNoUsername(t *testing.T) {
	tmpDir := t.TempDir()
	entry := crypt.CryptEntry{
		Salt:         "salt",
		Nonce:        "nonce",
		EncryptedDEK: "encrypted",
	}
	data, _ := json.Marshal(entry)

	err := crypt.ImportUser(tmpDir, string(data))
	if err == nil {
		t.Fatal("ImportUser without username should fail")
	}
}

func TestReadCryptFileNewFile(t *testing.T) {
	tmpDir := t.TempDir()
	cryptPath := filepath.Join(tmpDir, ".crypt")

	cf, err := crypt.ReadCryptFile(cryptPath)
	if err != nil {
		t.Fatalf("ReadCryptFile failed: %v", err)
	}

	if cf.Version != 1 {
		t.Errorf("Version mismatch: got %d, want 1", cf.Version)
	}

	if len(cf.Entries) != 0 {
		t.Errorf("Entries should be empty: got %d", len(cf.Entries))
	}
}

func TestReadCryptFileExisting(t *testing.T) {
	tmpDir := t.TempDir()
	setupTestUser(t, tmpDir, "user1", "pass1")

	cryptPath := filepath.Join(tmpDir, ".crypt")
	cf, err := crypt.ReadCryptFile(cryptPath)
	if err != nil {
		t.Fatalf("ReadCryptFile failed: %v", err)
	}

	if len(cf.Entries) != 1 {
		t.Errorf("Expected 1 entry, got %d", len(cf.Entries))
	}

	if cf.Entries[0].Username != "user1" {
		t.Errorf("Username mismatch: got %q, want %q", cf.Entries[0].Username, "user1")
	}
}

func TestWriteCryptFile(t *testing.T) {
	tmpDir := t.TempDir()
	cryptPath := filepath.Join(tmpDir, ".crypt")

	cf := crypt.CryptFile{
		Version: 1,
		Entries: []crypt.CryptEntry{
			{
				Username:     "testuser",
				Salt:         "salt",
				Nonce:        "nonce",
				EncryptedDEK: "encrypted",
			},
		},
	}

	err := crypt.WriteCryptFile(cryptPath, &cf)
	if err != nil {
		t.Fatalf("WriteCryptFile failed: %v", err)
	}

	// Verify file was created
	if _, err := os.Stat(cryptPath); err != nil {
		t.Fatalf("Crypt file not created: %v", err)
	}

	// Verify permissions
	info, _ := os.Stat(cryptPath)
	if info.Mode().Perm() != 0600 {
		t.Errorf("File permissions mismatch: got %o, want 600", info.Mode().Perm())
	}

	// Verify content
	cf2, err := crypt.ReadCryptFile(cryptPath)
	if err != nil {
		t.Fatalf("ReadCryptFile failed: %v", err)
	}

	if len(cf2.Entries) != 1 || cf2.Entries[0].Username != "testuser" {
		t.Fatal("Written data mismatch")
	}
}

func TestEncodeDecodeStorage(t *testing.T) {
	salt := []byte("sixteen bytes!!")
	nonce := []byte("twelve byte!")
	encrypted := []byte("some encrypted data here")

	// Encode
	saltB64, nonceB64, encryptedB64 := crypt.EncodeForStorage(salt, nonce, encrypted)

	// Verify base64
	if saltB64 == string(salt) {
		t.Fatal("Salt should be base64 encoded")
	}

	// Decode
	decodedSalt, decodedNonce, decodedEncrypted, err := crypt.DecodeFromStorage(saltB64, nonceB64, encryptedB64)
	if err != nil {
		t.Fatalf("DecodeFromStorage failed: %v", err)
	}

	// Verify roundtrip
	if string(decodedSalt) != string(salt) {
		t.Errorf("Salt mismatch: got %q, want %q", decodedSalt, salt)
	}

	if string(decodedNonce) != string(nonce) {
		t.Errorf("Nonce mismatch: got %q, want %q", decodedNonce, nonce)
	}

	if string(decodedEncrypted) != string(encrypted) {
		t.Errorf("Encrypted mismatch: got %q, want %q", decodedEncrypted, encrypted)
	}
}

func TestEncodeDecodeStorageInvalid(t *testing.T) {
	tests := []struct {
		name    string
		salt    string
		nonce   string
		encrypt string
	}{
		{"invalid salt", "!!!", "nonce", "encrypted"},
		{"invalid nonce", "salt", "!!!", "encrypted"},
		{"invalid encrypted", "salt", "nonce", "!!!"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, _, err := crypt.DecodeFromStorage(tt.salt, tt.nonce, tt.encrypt)
			if err == nil {
				t.Errorf("DecodeFromStorage should fail for %s", tt.name)
			}
		})
	}
}

func TestFindEntry(t *testing.T) {
	cf := crypt.CryptFile{
		Version: 1,
		Entries: []crypt.CryptEntry{
			{Username: "user1", Salt: "s1", Nonce: "n1", EncryptedDEK: "e1"},
			{Username: "user2", Salt: "s2", Nonce: "n2", EncryptedDEK: "e2"},
		},
	}

	// Find existing
	entry := cf.FindEntry("user1")
	if entry == nil {
		t.Fatal("FindEntry should find user1")
	}
	if entry.Username != "user1" {
		t.Errorf("Username mismatch: got %q, want user1", entry.Username)
	}

	// Find non-existent
	entry = cf.FindEntry("user3")
	if entry != nil {
		t.Fatal("FindEntry should return nil for non-existent user")
	}
}

func TestAddOrUpdateEntry(t *testing.T) {
	cf := crypt.CryptFile{
		Version: 1,
		Entries: []crypt.CryptEntry{
			{Username: "user1", Salt: "s1", Nonce: "n1", EncryptedDEK: "e1"},
		},
	}

	// Add new entry
	cf.AddOrUpdateEntry(crypt.CryptEntry{
		Username:     "user2",
		Salt:         "s2",
		Nonce:        "n2",
		EncryptedDEK: "e2",
	})

	if len(cf.Entries) != 2 {
		t.Errorf("Expected 2 entries, got %d", len(cf.Entries))
	}

	// Update existing entry
	cf.AddOrUpdateEntry(crypt.CryptEntry{
		Username:     "user1",
		Salt:         "s1_updated",
		Nonce:        "n1",
		EncryptedDEK: "e1",
	})

	if len(cf.Entries) != 2 {
		t.Errorf("Expected 2 entries after update, got %d", len(cf.Entries))
	}

	entry := cf.FindEntry("user1")
	if entry.Salt != "s1_updated" {
		t.Errorf("Salt not updated: got %q, want s1_updated", entry.Salt)
	}
}

func TestEncryptDeterminism(t *testing.T) {
	tmpDir := t.TempDir()
	setupTestUser(t, tmpDir, "testuser", "password")
	kp, _ := crypt.NewKeyProvider(tmpDir, "testuser", "password")

	plaintext := []byte("Test data")

	// Encrypt same data twice
	encrypted1, _ := kp.Encrypt(plaintext)
	encrypted2, _ := kp.Encrypt(plaintext)

	// Ciphertexts should be different (due to random nonce)
	if string(encrypted1) == string(encrypted2) {
		t.Fatal("Same plaintext should produce different ciphertexts (due to random nonce)")
	}

	// Both should decrypt to same plaintext
	decrypted1, _ := kp.Decrypt(encrypted1)
	decrypted2, _ := kp.Decrypt(encrypted2)

	if string(decrypted1) != string(plaintext) || string(decrypted2) != string(plaintext) {
		t.Fatal("Both ciphertexts should decrypt to same plaintext")
	}
}

func TestMultipleUsersIndependent(t *testing.T) {
	tmpDir := t.TempDir()
	setupTestUser(t, tmpDir, "user1", "pass1")
	setupTestUser(t, tmpDir, "user2", "pass2")

	kp1, _ := crypt.NewKeyProvider(tmpDir, "user1", "pass1")
	kp2, _ := crypt.NewKeyProvider(tmpDir, "user2", "pass2")

	plaintext := []byte("test")
	encrypted1, _ := kp1.Encrypt(plaintext)

	// User2 should not be able to decrypt user1's data
	_, err := kp2.Decrypt(encrypted1)
	if err == nil {
		t.Fatal("User2 should not decrypt user1's data")
	}

	// User1 should still be able to decrypt
	decrypted, err := kp1.Decrypt(encrypted1)
	if err != nil || string(decrypted) != string(plaintext) {
		t.Fatal("User1 should decrypt their own data")
	}
}
