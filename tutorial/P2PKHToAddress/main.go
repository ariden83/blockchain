package main

import (
	"crypto/sha256"
	"encoding/hex"
	"log"

	"github.com/btcsuite/btcutil/base58"
)

// P2PKHToAddress get the bitcoin address from the Pay To Public Key Hash script
// https://bitcoin.stackexchange.com/questions/19081/parsing-bitcoin-input-and-output-addresses-from-scripts
func P2PKHToAddress(pkscript []byte, isTestnet bool) (string, error) {
	p := make([]byte, 1)
	p[0] = 0x80 // prefix with 00 if it's mainnet
	if isTestnet {
		p[0] = 0x6F // prefix with 0F if it's testnet
	}
	pf := append(p[:], pkscript[:]...) // add prefix
	h1 := sha256.Sum256(pf)            // hash it
	h2 := sha256.Sum256(h1[:])         // hash it again
	b := append(pf[:], h2[0:4]...)     // prepend the prefix to the first 5 bytes
	address := base58.Encode(b)        // encode to base58
	if !isTestnet {
		address = "1" + address // prefix with 1 if it's mainnet
	}

	return address, nil
}

func main() {
	hash, err := hex.DecodeString("76a914877fefc337afdc98afe4a4ab1c4e85221292783988ac")
	if err != nil {
		log.Fatal(err)
	}
	address, err := P2PKHToAddress(hash, true)
	if err != nil {
		log.Fatal(err)
	}
	// testnet address
	expected := "mssQn95JGBtw6Npbt6Z8LoJfK1Buuz6ZHt"
	if address != expected {
		log.Fatalf("got: %v; expected: %v", address, expected)
	}

	log.Printf("address: %v", address)
}
