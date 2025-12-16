package keystore

import "errors"

const (
	KeychainService = "com.pastee.clipboard"
	KeychainAccount = "database-encryption-key"
)

var (
	ErrKeyNotFound    = errors.New("encryption key not found in keychain")
	ErrKeyStoreFailed = errors.New("failed to store key in keychain")
)

type KeyStore interface {
	Get() ([]byte, error)
	Set(key []byte) error
	Delete() error
	Exists() (bool, error)
}

func NewKeyStore() KeyStore {
	return newPlatformKeyStore()
}
