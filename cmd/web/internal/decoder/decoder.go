package decoder

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
)

var privateKey = ""

func GetPrivateKey() string {
	if privateKey != "" {
		return privateKey
	}
	key := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		return GetPrivateKey()
	}

	privateKey = base64.StdEncoding.EncodeToString(key)
	return privateKey
}

func Encrypt(plaintext []byte, key64 string) (string, error) {

	key, err := base64.StdEncoding.DecodeString(key64)
	if err != nil {
		return "", err
	}

	c, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func Decrypt(cipherTxt64 string, key64 string) ([]byte, error) {

	ciphertext, err := base64.StdEncoding.DecodeString(cipherTxt64)
	if err != nil {
		return []byte{}, err
	}
	key, err := base64.StdEncoding.DecodeString(key64)
	if err != nil {
		return []byte{}, err
	}

	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, nil)
}
