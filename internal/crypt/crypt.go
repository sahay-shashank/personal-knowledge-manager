package crypt

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// ReadCryptFile reads .crypt file or returns empty CryptFile if not exists
func ReadCryptFile(cryptPath string) (*CryptFile, error) {
	
	data, err := os.ReadFile(cryptPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			// New crypt file
			return &CryptFile{Version: 1, Entries: []CryptEntry{}}, nil
		}
		return nil, err
	}

	var cf CryptFile
	if err := json.Unmarshal(data, &cf); err != nil {
		return nil, fmt.Errorf("corrupt .crypt file: %w", err)
	}
	return &cf, nil
}

// WriteCryptFile writes CryptFile to disk
func WriteCryptFile(cryptPath string, cf *CryptFile) error {
	// Create directory if not exists
	if err := os.MkdirAll(filepath.Dir(cryptPath), 0700); err != nil {
		return err
	}

	data, err := json.MarshalIndent(cf, "", "  ")
	if err != nil {
		return err
	}

	// Write with restricted permissions (user only)
	return os.WriteFile(cryptPath, data, 0600)
}

// FindEntry finds user entry in CryptFile
func (cf *CryptFile) FindEntry(username string) *CryptEntry {
	for i := range cf.Entries {
		if cf.Entries[i].Username == username {
			return &cf.Entries[i]
		}
	}
	return nil
}

// AddOrUpdateEntry adds or updates user entry
func (cf *CryptFile) AddOrUpdateEntry(entry CryptEntry) {
	for i := range cf.Entries {
		if cf.Entries[i].Username == entry.Username {
			cf.Entries[i] = entry
			return
		}
	}
	cf.Entries = append(cf.Entries, entry)
}

// EncodeForStorage encodes salt/nonce/encrypted values to base64
func EncodeForStorage(salt, nonce, encrypted []byte) (saltB64, nonceB64, encryptedB64 string) {
	return base64.StdEncoding.EncodeToString(salt),
		base64.StdEncoding.EncodeToString(nonce),
		base64.StdEncoding.EncodeToString(encrypted)
}

// DecodeFromStorage decodes base64 values back to bytes
func DecodeFromStorage(saltB64, nonceB64, encryptedB64 string) (salt, nonce, encrypted []byte, err error) {
	salt, err = base64.StdEncoding.DecodeString(saltB64)
	if err != nil {
		return nil, nil, nil, err
	}

	nonce, err = base64.StdEncoding.DecodeString(nonceB64)
	if err != nil {
		return nil, nil, nil, err
	}

	encrypted, err = base64.StdEncoding.DecodeString(encryptedB64)
	if err != nil {
		return nil, nil, nil, err
	}

	return salt, nonce, encrypted, nil
}
