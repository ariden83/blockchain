package main

import (
	"crypto/aes"
	"fmt"
)

// TODO: support n byte pad (see TLS), not only 0 byte pad
// TODO: (DES): Support different key length and key derivations, e.g. EDE2 etc.
// TODO: check key lengths at all (16 bytes aes, 8 bytes des, 3*8 bytes 3des)
// TODO: write tests

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

func main() {
	// Some testing

	//key := []byte("1234567890123456") // aes, 16 bytes
	key := []byte("1234567812345678f1235678") // des, 8 bytes
	plaintext := []byte("Freda is the name of a cow.")
	//iv := generateIV(aes.BlockSize)
	iv := []byte("1234567890123456")
	fmt.Printf("\nKey: %v\nPlaintext: %v\nIV: %0x\n\n", string(key), string(plaintext), string(iv))

	ciphertext, err := Encrypt(plaintext, key, nil, TDES, 0)
	if err != nil {
		fmt.Errorf("Couldn't encrypt: %v", err.Error())
	}
	fmt.Printf("Ciphertext: %v\n\n", ciphertext)

	plaintext, err = Decrypt(ciphertext, key, TDES, 0)
	if err != nil {
		fmt.Errorf("Couldn't decrypt: %v", err.Error())
	}
	fmt.Printf("Decryption result: %v\n\n", string(plaintext))
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
