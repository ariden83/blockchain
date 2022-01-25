package decoder

import (
	"crypto/aes"
	"crypto/cipher"
	"errors"
	"crypto/des"
	"encoding/base64"
)

type Cipher int
type Mode int

// Assign uints to different ciphers and modes of operation
const (
	AES Cipher = iota
	DES
	TDES

	CBC Mode = iota
	CTR
	CFB
	OFB
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
func Encrypt(plaintext, key, iv []byte, cipher Cipher, mode Mode) ([]byte, error) {
	input := []byte(plaintext)
	var output []byte
	var err error

	switch cipher {
	case AES:
		output = make([]byte, aes.BlockSize+len(input))
		err = encryptAES(input, output, key, iv, mode)
	case DES:
		// No need to pass output slice because its length
		// depends on input slice length and its padding.
		// Thus output is created after appending the padding.
		output, err = encryptDES(input, key, iv, false)
	case TDES:
		output, err = encryptDES(input, key, iv, true) // last parameter indicates use of 3DES
	}

	if err != nil {
		return nil, err
	}
	return output, nil
}

// Decrypt serves as wrapper function for decrypting any ciphertext,key with specified
// cipher and mode of operation
func Decrypt(ciphertext, key []byte, cipher Cipher, mode Mode) ([]byte, error) {
	input := []byte(ciphertext)
	var output []byte
	var err error

	switch cipher {
	case AES:
		output = make([]byte, len(input)-aes.BlockSize)
		err = decryptAES(input, output, key, mode)
	case DES:
		// Here we can create output before decryption
		// because its length is the same as input's length
		output = make([]byte, len(input))
		err = decryptDES(input, output, key, false)
	case TDES:
		output = make([]byte, len(input))
		err = decryptDES(input, output, key, true) // last parameter indicates use of 3DES
	}

	if err != nil {
		return nil, err
	}
	return output, nil
}


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


// encryptDES enrypts plaintext input with passed key, IV and 3DES-flag in DES block cipher;
// and returns ciphertext output
func encryptDES(input, key, iv []byte, tripleDES bool) ([]byte, error) {
	var block cipher.Block
	var err error
	if tripleDES {
		block, err = des.NewTripleDESCipher(key)
	} else {
		block, err = des.NewCipher(key)
	}

	if err != nil {
		return nil, errors.New("Couldn't create block cipher.")
	}
	if len(input)%des.BlockSize != 0 {
		input = addPadding(input, des.BlockSize)
	}
	output := make([]byte, len(input))

	for i := 0; i < len(input)/des.BlockSize; i++ {
		start := des.BlockSize * i
		end := start + des.BlockSize
		block.Encrypt(output[start:end], input[start:end])
	}
	return output, nil
}

// decryptDES derypts ciphertext input with passed key and 3DES-flag
// in DES block cipher; and returns plaintext output
func decryptDES(input, output, key []byte, tripleDES bool) error {
	var block cipher.Block
	var err error
	if tripleDES {
		block, err = des.NewTripleDESCipher(key)
	} else {
		block, err = des.NewCipher(key)
	}

	if err != nil {
		return errors.New("Couldn't create block cipher.")
	}

	for i := 0; i < len(input)/des.BlockSize; i++ {
		start := des.BlockSize * i
		end := start + des.BlockSize
		block.Decrypt(output[start:end], input[start:end])
	}
	return nil
}
