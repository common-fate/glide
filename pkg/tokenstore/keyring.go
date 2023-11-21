package tokenstore

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/99designs/keyring"
	"github.com/common-fate/clio"

	"github.com/pkg/errors"
)

// cfKeyring is a wrapper around 99designs/keyring
// that handles config via env vars and
// marshalling/unmarshalling of items.
type cfKeyring struct {
	// keyring is an existing keyring to use.
	// if nil, a new keyring is created using the openKeyring() method.
	keyring keyring.Keyring
}

// returns false if the key is not found, true if it is found, or false and an error if there was a keyring related error
func (s *cfKeyring) HasKey(key string) (bool, error) {
	ring, err := s.openKeyring()
	if err != nil {
		return false, err
	}
	_, err = ring.Get(key)
	if err == keyring.ErrKeyNotFound {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

// returns keyring.ErrKeyNotFound if not found
func (s *cfKeyring) Retrieve(key string, target interface{}) error {
	ring, err := s.openKeyring()
	if err != nil {
		return err
	}
	keyringItem, err := ring.Get(key)
	if err != nil {
		return err
	}
	return json.Unmarshal(keyringItem.Data, &target)
}

func (s *cfKeyring) Store(key string, payload interface{}) error {
	ring, err := s.openKeyring()
	if err != nil {
		return err
	}
	b, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	return ring.Set(keyring.Item{
		Key:  key, // store with the corresponding key
		Data: b,   // store the bytes
	})
}

func (s *cfKeyring) Clear(key string) error {
	ring, err := s.openKeyring()
	if err != nil {
		return err
	}
	return ring.Remove(key)
}

func (s *cfKeyring) List() ([]keyring.Item, error) {
	tokenList := []keyring.Item{}
	ring, err := s.openKeyring()
	if err != nil {
		return nil, err
	}
	keys, err := ring.Keys()
	if err != nil {
		return nil, err
	}
	for _, k := range keys {
		item, err := ring.Get(k)
		if err != nil {
			return nil, err
		}
		tokenList = append(tokenList, item)

	}
	return tokenList, nil
}

func (s *cfKeyring) ListKeys() ([]string, error) {
	ring, err := s.openKeyring()
	if err != nil {
		return nil, err
	}
	return ring.Keys()
}

func (s *cfKeyring) openKeyring() (keyring.Keyring, error) {
	// return the existing keyring if there's one set.
	if s.keyring != nil {
		clio.Debug("existing keyring has been set: returning")
		return s.keyring, nil
	}

	name := os.Getenv("COMMONFATE_KEYRING_NAME")
	if name == "" {
		name = "commonfate"
	}

	dirname := os.Getenv("COMMONFATE_KEYRING_FILE_DIR")
	if dirname == "" {
		dirname = "~/.commonfate"
	}

	keychainName := os.Getenv("COMMONFATE_KEYRING_MACOS_KEYCHAIN_NAME")
	if keychainName == "" {
		keychainName = "login"
	}

	c := keyring.Config{
		ServiceName: name,

		// MacOS keychain
		KeychainName:             keychainName,
		KeychainTrustApplication: true,

		// KDE Wallet
		KWalletAppID:  name,
		KWalletFolder: name,

		// Windows
		WinCredPrefix: name,

		// freedesktop.org's Secret Service
		LibSecretCollectionName: name,

		// Pass (https://www.passwordstore.org/)
		PassPrefix: name,

		// Fallback encrypted file
		FileDir:          dirname,
		FilePasswordFunc: keyring.TerminalPrompt,
	}

	kab := os.Getenv("COMMONFATE_KEYRING_ALLOWED_BACKENDS")
	if kab != "" {
		kab = strings.ReplaceAll(kab, " ", "") // remove any spaces
		backends := strings.Split(kab, ",")

		clio.Debugw("setting allowed keyring backends", "backends", backends)

		for _, b := range backends {
			c.AllowedBackends = append(c.AllowedBackends, keyring.BackendType(b))
		}
	} else {
		backends := keyring.AvailableBackends()
		clio.Debugw("setting default keyring backends", "backends", backends)
		c.AllowedBackends = backends
	}

	if strings.ToLower(os.Getenv("COMMONFATE_KEYRING_DEBUG")) == "true" {
		keyring.Debug = true
	}

	k, err := keyring.Open(c)
	if err != nil {
		return nil, errors.Wrap(err, "opening keyring")
	}

	return k, nil
}
