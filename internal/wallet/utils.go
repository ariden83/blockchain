package wallet

import (
	"errors"

	"golang.org/x/crypto/bcrypt"

	"github.com/wemeetagain/go-hdwallet"
	"golang.org/x/crypto/sha3"
)

func GetPubKey(privKey []byte) ([]byte, error) {
	masterPrv, err := hdwallet.StringWallet(string(privKey))
	if err != nil {
		return []byte{}, err
	}
	return []byte(masterPrv.Pub().String()), nil
}

func encryptPassword(password []byte) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword(password, bcrypt.MinCost)
	if err != nil {
		return []byte{}, errors.New("fail to encrypt password")
	}
	return hash, nil
}

const localKey = "blockchain"

func hash(mnemonic []byte) []byte {
	fixedSlice := sha3.Sum512(append(mnemonic, []byte(localKey)...))
	byteSlice := make([]byte, 64)
	for i := 0; i < len(fixedSlice); i++ {
		byteSlice[i] = fixedSlice[i]
	}
	return byteSlice
}
