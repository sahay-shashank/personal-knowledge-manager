package crypt

const (
	DEKSize   = 32 // AES-256
	KEKSize   = 32 // AES-256
	SaltSize  = 16
	PBKDFIter = 100000
)

type KeyProvider struct {
	username string
	dek []byte // in-memory DEK

}

type CryptEntry struct {
	Username     string `json:"username"`
	Salt         string `json:"salt"`          // base64 encoded
	Nonce        string `json:"nonce"`         // base64 encoded
	EncryptedDEK string `json:"encrypted_dek"` // base64 encoded
}

type CryptFile struct {
	Version int          `json:"version"`
	Entries []CryptEntry `json:"entries"`
}
