package main

import (
	"crypto/aes"
	"crypto/cipher"
	"errors"
)

// encryptAES enrypts plaintext input with passed key, IV and mode in AES block cipher;
// and returns ciphertext output
func encryptAES(input []byte, output []byte, key, iv []byte, mode Mode) error {
	block, err := aes.NewCipher(key)
	if err != nil {
		return errors.New("Couldn't create block cipher.")
	}

	// Prepend IV to ciphertext.
	// Generate IV randomly if it is not passed
	if iv == nil {
		if iv = generateIV(aes.BlockSize); iv == nil {
			return errors.New("Couldn't create random initialization vector (IV).")
		}
	}
	copy(output, iv)

	switch mode {
	case CBC:
		if len(input)%aes.BlockSize != 0 {
			input = addPadding(input, aes.BlockSize)
		}
		mode := cipher.NewCBCEncrypter(block, iv)
		mode.CryptBlocks(output[aes.BlockSize:], input)
	case CFB:
		mode := cipher.NewCFBEncrypter(block, iv)
		mode.XORKeyStream(output[aes.BlockSize:], input)
	case CTR:
		mode := cipher.NewCTR(block, iv)
		mode.XORKeyStream(output[aes.BlockSize:], input)
	case OFB:
		mode := cipher.NewOFB(block, iv)
		mode.XORKeyStream(output[aes.BlockSize:], input)
	}
	return nil
}

// decryptAES derypts ciphertext input with passed key and mode (IV is contained in input)
// in AES block cipher; and returns plaintext output
func decryptAES(input []byte, output []byte, key []byte, mode Mode) error {
	block, err := aes.NewCipher(key)
	if err != nil {
		return errors.New("Couldn't create block cipher.")
	}
	if len(input) < aes.BlockSize {
		return errors.New("Ciphertext too short.")
	}

	iv := input[:aes.BlockSize]
	ciphertext := input[aes.BlockSize:]

	switch mode {
	case CBC:
		if len(input)%aes.BlockSize != 0 {
			return errors.New("Ciphertext doesn't satisfy CBC-mode requirements.")
		}
		mode := cipher.NewCBCDecrypter(block, iv)
		mode.CryptBlocks(output, ciphertext)
	case CFB:
		mode := cipher.NewCFBDecrypter(block, iv)
		mode.XORKeyStream(output, ciphertext)
	case CTR:
		mode := cipher.NewCTR(block, iv)
		mode.XORKeyStream(output, ciphertext)
	case OFB:
		mode := cipher.NewOFB(block, iv)
		mode.XORKeyStream(output, ciphertext)
	}
	return nil
}
