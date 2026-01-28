package crypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/crypto/pbkdf2"
)

// NewKeyProvider returns a provider
func NewKeyProvider(pkmDir string, username string, password string) (*KeyProvider, error) {
	if username == "" {
		return nil, fmt.Errorf("username required")
	}

	cryptPath := filepath.Join(pkmDir, ".crypt")

	// Read .crypt file
	cf, err := ReadCryptFile(cryptPath)
	if err != nil {
		return nil, err
	}

	// Find user entry
	entry := cf.FindEntry(username)
	if entry == nil {
		return nil, fmt.Errorf("user %q not found in .crypt. Run 'ztl init --user %s' first", username, username)
	}

	// Decode stored values
	salt, nonce, encryptedDEK, err := DecodeFromStorage(entry.Salt, entry.Nonce, entry.EncryptedDEK)
	if err != nil {
		return nil, fmt.Errorf("corrupt .crypt entry: %w", err)
	}

	// Derive KEK from password + salt
	kek := pbkdf2.Key([]byte(password), salt, PBKDFIter, KEKSize, sha256.New)

	// Decrypt DEK
	dek, err := decryptAESGCM(kek, nonce, encryptedDEK)
	if err != nil {
		return nil, fmt.Errorf("decrypt DEK failed (wrong password?): %w", err)
	}

	return &KeyProvider{
		username: username,
		dek:      dek,
	}, nil

}
func InitUser(pkmDir, username, password string) error {
	if username == "" || password == "" {
		return fmt.Errorf("username and password required")
	}

	cryptPath := filepath.Join(pkmDir, ".crypt")

	// Read existing .crypt
	cf, err := ReadCryptFile(cryptPath)
	if err != nil {
		return err
	}

	// Check if user already exists
	if cf.FindEntry(username) != nil {
		return fmt.Errorf("user %q already exists", username)
	}

	// Generate random salt
	salt := make([]byte, SaltSize)
	if _, err := rand.Read(salt); err != nil {
		return err
	}

	// Derive KEK from password + salt
	kek := pbkdf2.Key([]byte(password), salt, PBKDFIter, KEKSize, sha256.New)

	// Generate random DEK
	dek := make([]byte, DEKSize)
	if _, err := rand.Read(dek); err != nil {
		return err
	}

	// Encrypt DEK with KEK
	nonce := make([]byte, 12)
	if _, err := rand.Read(nonce); err != nil {
		return err
	}

	encryptedDEK, err := encryptAESGCM(kek, nonce, dek)
	if err != nil {
		return err
	}

	// Encode for storage
	saltB64, nonceB64, encryptedB64 := EncodeForStorage(salt, nonce, encryptedDEK)

	// Add entry to .crypt
	entry := CryptEntry{
		Username:     username,
		Salt:         saltB64,
		Nonce:        nonceB64,
		EncryptedDEK: encryptedB64,
	}
	cf.AddOrUpdateEntry(entry)

	// Write .crypt
	if err := WriteCryptFile(cryptPath, cf); err != nil {
		return err
	}

	// Create user directory
	userDir := filepath.Join(pkmDir, username)
	if err := os.MkdirAll(userDir, 0755); err != nil {
		return err
	}

	return nil
}

// ChangePassword updates user's password
func ChangePassword(pkmDir, username, oldPassword, newPassword string) error {
	if newPassword == "" {
		return fmt.Errorf("new password required")
	}

	cryptPath := filepath.Join(pkmDir, ".crypt")

	// Load with old password
	kp, err := NewKeyProvider(pkmDir, username, oldPassword)
	if err != nil {
		return fmt.Errorf("old password incorrect: %w", err)
	}

	// Read .crypt
	cf, err := ReadCryptFile(cryptPath)
	if err != nil {
		return err
	}

	entry := cf.FindEntry(username)
	if entry == nil {
		return fmt.Errorf("user not found")
	}

	// Generate new salt
	newSalt := make([]byte, SaltSize)
	if _, err := rand.Read(newSalt); err != nil {
		return err
	}

	// Derive new KEK with new password
	newKEK := pbkdf2.Key([]byte(newPassword), newSalt, PBKDFIter, KEKSize, sha256.New)

	// Encrypt same DEK with new KEK
	newNonce := make([]byte, 12)
	if _, err := rand.Read(newNonce); err != nil {
		return err
	}

	newEncryptedDEK, err := encryptAESGCM(newKEK, newNonce, kp.dek)
	if err != nil {
		return err
	}

	// Update entry
	saltB64, nonceB64, encryptedB64 := EncodeForStorage(newSalt, newNonce, newEncryptedDEK)
	entry.Salt = saltB64
	entry.Nonce = nonceB64
	entry.EncryptedDEK = encryptedB64

	// Write back
	return WriteCryptFile(cryptPath, cf)
}

// DEK returns the data encryption key
func (kp *KeyProvider) DEK() []byte {
	return kp.dek
}

// Encrypt encrypts plaintext with session DEK
func (kp *KeyProvider) Encrypt(plaintext []byte) ([]byte, error) {
	nonce := make([]byte, 12)
	if _, err := rand.Read(nonce); err != nil {
		return nil, err
	}

	ciphertext, err := encryptAESGCM(kp.dek, nonce, plaintext)
	if err != nil {
		return nil, err
	}

	return append(nonce, ciphertext...), nil
}

// Decrypt decrypts ciphertext with session DEK
func (kp *KeyProvider) Decrypt(ciphertext []byte) ([]byte, error) {
	if len(ciphertext) < 12 {
		return nil, fmt.Errorf("ciphertext too short")
	}

	nonce := ciphertext[:12]
	encrypted := ciphertext[12:]

	return decryptAESGCM(kp.dek, nonce, encrypted)
}

// Helper: encrypt with AES-GCM
func encryptAESGCM(key, nonce, plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	return aesgcm.Seal(nil, nonce, plaintext, nil), nil
}

// Helper: decrypt with AES-GCM
func decryptAESGCM(key, nonce, ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	return aesgcm.Open(nil, nonce, ciphertext, nil)
}

func (kp *KeyProvider) GetDEK() []byte {
	return kp.dek
}

// Export user profile from the current host
func ExportUser(pkmDir, username, password string) error {
	if password == "" {
		return fmt.Errorf("Password required")
	}

	cryptPath := filepath.Join(pkmDir, ".crypt")

	// Read .crypt
	cf, err := ReadCryptFile(cryptPath)
	if err != nil {
		return err
	}

	entry := cf.FindEntry(username)
	if entry == nil {
		return fmt.Errorf("user not found")
	}

	// Validate user with password
	if _, err := NewKeyProvider(pkmDir, username, password); err != nil {
		return fmt.Errorf("Password incorrect: %w", err)
	}

	output, err := json.Marshal(entry)
	fmt.Println(string(output))
	return nil
}

// Import user profile to the current host
func ImportUser(pkmDir, userData string) error {

	entry := CryptEntry{}
	if err := json.Unmarshal([]byte(userData), &entry); err != nil {
		return err
	}

	if entry.Username == "" {
		return fmt.Errorf("username not found in input")
	}

	cryptPath := filepath.Join(pkmDir, ".crypt")

	// Read existing .crypt
	cf, err := ReadCryptFile(cryptPath)
	if err != nil {
		return err
	}

	// Check if user already exists
	if cf.FindEntry(entry.Username) != nil {
		return fmt.Errorf("user %q already exists", entry.Username)
	}
	cf.AddOrUpdateEntry(entry)

	// Write .crypt
	if err := WriteCryptFile(cryptPath, cf); err != nil {
		return err
	}

	// Create user directory
	userDir := filepath.Join(pkmDir, entry.Username)
	if err := os.MkdirAll(userDir, 0755); err != nil {
		return err
	}

	return nil
}
