//go:build linux

package keystore

import (
	"github.com/zalando/go-keyring"
)

type linuxKeyStore struct{}

func newPlatformKeyStore() KeyStore {
	return &linuxKeyStore{}
}

func (k *linuxKeyStore) Get() ([]byte, error) {
	secret, err := keyring.Get(KeychainService, KeychainAccount)
	if err == keyring.ErrNotFound {
		return nil, ErrKeyNotFound
	}
	if err != nil {
		return nil, err
	}
	return []byte(secret), nil
}

func (k *linuxKeyStore) Set(key []byte) error {
	err := keyring.Set(KeychainService, KeychainAccount, string(key))
	if err != nil {
		return ErrKeyStoreFailed
	}
	return nil
}

func (k *linuxKeyStore) Delete() error {
	return keyring.Delete(KeychainService, KeychainAccount)
}

func (k *linuxKeyStore) Exists() (bool, error) {
	_, err := keyring.Get(KeychainService, KeychainAccount)
	if err == keyring.ErrNotFound {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}
