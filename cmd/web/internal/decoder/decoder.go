package decoder

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
)

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

	println(string(decodedText))
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
