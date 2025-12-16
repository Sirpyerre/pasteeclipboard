package keystore

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

func GenerateEncryptionKey() (string, error) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		return "", fmt.Errorf("failed to generate random key: %w", err)
	}

	return hex.EncodeToString(key), nil
}

func GetOrCreateKey(store KeyStore) (string, error) {
	exists, err := store.Exists()
	if err != nil {
		return "", fmt.Errorf("failed to check key existence: %w", err)
	}

	if exists {
		keyBytes, err := store.Get()
		if err != nil {
			return "", fmt.Errorf("failed to retrieve existing key: %w", err)
		}
		return string(keyBytes), nil
	}

	key, err := GenerateEncryptionKey()
	if err != nil {
		return "", err
	}

	if err := store.Set([]byte(key)); err != nil {
		return "", fmt.Errorf("failed to store key in keychain: %w", err)
	}

	return key, nil
}
