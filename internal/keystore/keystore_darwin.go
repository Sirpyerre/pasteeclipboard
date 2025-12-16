//go:build darwin

package keystore

import (
	"github.com/keybase/go-keychain"
)

type darwinKeyStore struct{}

func newPlatformKeyStore() KeyStore {
	return &darwinKeyStore{}
}

func (k *darwinKeyStore) Get() ([]byte, error) {
	query := keychain.NewItem()
	query.SetSecClass(keychain.SecClassGenericPassword)
	query.SetService(KeychainService)
	query.SetAccount(KeychainAccount)
	query.SetMatchLimit(keychain.MatchLimitOne)
	query.SetReturnData(true)

	results, err := keychain.QueryItem(query)
	if err == keychain.ErrorItemNotFound {
		return nil, ErrKeyNotFound
	}
	if err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return nil, ErrKeyNotFound
	}

	return results[0].Data, nil
}

func (k *darwinKeyStore) Set(key []byte) error {
	k.Delete()

	item := keychain.NewItem()
	item.SetSecClass(keychain.SecClassGenericPassword)
	item.SetService(KeychainService)
	item.SetAccount(KeychainAccount)
	item.SetLabel("Pastee Clipboard Encryption Key")
	item.SetData(key)
	item.SetSynchronizable(keychain.SynchronizableNo)
	item.SetAccessible(keychain.AccessibleWhenUnlocked)

	err := keychain.AddItem(item)
	if err != nil {
		return ErrKeyStoreFailed
	}

	return nil
}

func (k *darwinKeyStore) Delete() error {
	item := keychain.NewItem()
	item.SetSecClass(keychain.SecClassGenericPassword)
	item.SetService(KeychainService)
	item.SetAccount(KeychainAccount)

	return keychain.DeleteItem(item)
}

func (k *darwinKeyStore) Exists() (bool, error) {
	_, err := k.Get()
	if err == ErrKeyNotFound {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}
