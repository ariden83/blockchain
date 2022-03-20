package signschnorr

import (
	"crypto/sha256"

	"math/big"
	"os"
	"testing"

	"github.com/LuisAcerv/btchdwallet/crypt"
	"github.com/brianium/mnemonic"
	"github.com/gcash/bchd/bchec"
	"github.com/stretchr/testify/assert"
	"github.com/wemeetagain/go-hdwallet"
)

func TestMain(m *testing.M) {
	exitVal := m.Run()
	os.Exit(exitVal)
}

// This example demonstrates creating a script which pays to a bitcoin address.
// It also prints the created script hex and uses the DisasmString function to
// display the disassembled script.
func Test_Example_(t *testing.T) {
	var (
		sign    *Signature
		pubkey  *bchec.PublicKey
		privKey *bchec.PrivateKey
		hash    []byte
	)

	// @https://github.com/gcash/bchd/blob/87534217bfc8d2c73461f025b5330f145aed9e86/bchec/signature_test.go
	t.Run("example SignSchnorr", func(t *testing.T) {
		seed := crypt.CreateHash()

		mnemonic, err := mnemonic.New([]byte(seed), mnemonic.English)
		assert.NoError(t, err)

		// Create a master private key
		masterPrv := hdwallet.MasterKey([]byte(mnemonic.Sentence()))
		privKey, pubkey = bchec.PrivKeyFromBytes(bchec.S256(), masterPrv.Key)

		message := "Satoshi Nakamoto"
		h := sha256.Sum256([]byte(message))
		hash = h[:]

		//sign, err = privKey.SignSchnorr(hash)
		sign, err = signSchnorr(privKey, hash)
		assert.NoError(t, err)
	})

	t.Run("example compress and decompress pubkey", func(t *testing.T) {
		compressKey := pubkey.SerializeCompressed()

		pk, err := bchec.ParsePubKey(compressKey, bchec.S256())
		assert.NoError(t, err)

		pubkey = pk
	})

	t.Run("example SignSchnorr verify", func(t *testing.T) {
		sig, err := bchec.ParseSchnorrSignature(sign.Serialize())
		assert.NoError(t, err)
		valid := sig.Verify(hash, pubkey)
		assert.True(t, valid)
	})
}

// signSchnorr signs the hash using the schnorr signature algorithm.
func signSchnorr(privateKey *bchec.PrivateKey, hash []byte) (*Signature, error) {
	// The rfc6979 nonce derivation function accepts additional entropy.
	// We are using the same entropy that is used by bitcoin-abc so our test
	// vectors will be compatible. This byte string is chosen to avoid collisions
	// with ECDSA which would render the signature insecure.
	//
	// See https://github.com/bitcoincashorg/bitcoincash.org/blob/master/spec/2019-05-15-schnorr.md#recommended-practices-for-secure-signature-generation
	additionalData := []byte{'S', 'c', 'h', 'n', 'o', 'r', 'r', '+', 'S', 'H', 'A', '2', '5', '6', ' ', ' '}
	k := nonceRFC6979(privateKey.D, hash, additionalData)
	// Compute point R = k * G
	rx, ry := privateKey.Curve.ScalarBaseMult(k.Bytes())

	//  Negate nonce if R.y is not a quadratic residue.
	if big.Jacobi(ry, privateKey.Params().P) != 1 {
		k = k.Neg(k)
	}

	// Compute scalar e = Hash(R.x || compressed(P) || m) mod N
	eBytes := sha256.Sum256(append(append(padIntBytes(rx), privateKey.PubKey().SerializeCompressed()...), hash...))
	e := new(big.Int).SetBytes(eBytes[:])
	e.Mod(e, privateKey.Params().N)

	// Compute scalar s = (k + e * x) mod N
	x := new(big.Int).SetBytes(privateKey.Serialize())
	s := e.Mul(e, x)
	s.Add(s, k)
	s.Mod(s, privateKey.Params().N)

	return &Signature{
		R: rx,
		S: s,
	}, nil
}
