package decoder

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
)

type Cipher int
type Mode int

func Password(ciphertext, iv, key string) (string, error) {
	decodedText, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}
	decodedIv, err := base64.StdEncoding.DecodeString(iv)
	if err != nil {
		return "", err
	}

	newCipher, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}
	cfbdec := cipher.NewCBCDecrypter(newCipher, decodedIv)
	cfbdec.CryptBlocks(decodedText, decodedText)

	decodedText = removeBadPadding(decodedText)

	//println(string(decodedText))
	data, err := base64.RawStdEncoding.DecodeString(string(decodedText))
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func removeBadPadding(b64 []byte) []byte {
	last := b64[len(b64)-1]
	if last > 16 {
		return b64
	}
	return b64[:len(b64)-int(last)]
}

// Encrypt serves as wrapper function for encrypting any plaintext,key with specified
// cipher and mode of operation
func Encrypt(plaintext, key []byte) ([]byte, []byte, error) {
	if len(plaintext)%aes.BlockSize != 0 {
		plaintext = addPadding(plaintext, aes.BlockSize)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, nil, err
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, nil, err
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], plaintext)

	return ciphertext, iv, err
}
