package wallet

import (
	"golang.org/x/crypto/bcrypt"

	pkgError "github.com/ariden83/blockchain/pkg/errors"
)

type SeedNoPrivKey struct {
	Timestamp int64
	Address   []byte
	Password  []byte
}

// Seed represents each 'item' in the blockchain
type Seed struct {
	Address   []byte
	Timestamp int64
	PubKey    []byte
	PrivKey   []byte
	Mnemonic  []byte
	ExtraData map[string]interface{}
}

func (s *Seed) verifyPassword(password []byte) error {
	if _, ok := s.ExtraData["password"]; !ok {
		return pkgError.ErrorSeedPasswordInvalid
	}

	if err := bcrypt.CompareHashAndPassword(s.ExtraData["password"].([]byte), password); err != nil {
		return pkgError.ErrorSeedPasswordInvalid
	}
	return nil
}
