//go:build windows

package keystore

import (
	"github.com/danieljoos/wincred"
)

type windowsKeyStore struct{}

func newPlatformKeyStore() KeyStore {
	return &windowsKeyStore{}
}

func (k *windowsKeyStore) Get() ([]byte, error) {
	cred, err := wincred.GetGenericCredential(KeychainService)
	if err != nil {
		return nil, ErrKeyNotFound
	}
	return cred.CredentialBlob, nil
}

func (k *windowsKeyStore) Set(key []byte) error {
	cred := wincred.NewGenericCredential(KeychainService)
	cred.UserName = KeychainAccount
	cred.CredentialBlob = key
	cred.Comment = "Pastee Clipboard Database Encryption Key"

	err := cred.Write()
	if err != nil {
		return ErrKeyStoreFailed
	}
	return nil
}

func (k *windowsKeyStore) Delete() error {
	cred, err := wincred.GetGenericCredential(KeychainService)
	if err != nil {
		return nil
	}
	return cred.Delete()
}

func (k *windowsKeyStore) Exists() (bool, error) {
	_, err := wincred.GetGenericCredential(KeychainService)
	if err != nil {
		return false, nil
	}
	return true, nil
}
