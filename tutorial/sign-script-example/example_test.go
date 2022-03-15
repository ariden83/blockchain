// Copyright (c) 2014-2016 The btcsuite developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package txscript_test

import (
	"bytes"
	"crypto/elliptic"
	"crypto/hmac"
	"crypto/sha256"
	"github.com/stretchr/testify/assert"

	"fmt"
	"hash"
	"math/big"

	"encoding/hex"
	"github.com/LuisAcerv/btchdwallet/crypt"
	"github.com/brianium/mnemonic"
	"github.com/gcash/bchd/bchec"
	"github.com/gcash/bchd/chaincfg"
	"github.com/gcash/bchd/chaincfg/chainhash"
	"github.com/gcash/bchd/txscript"
	"github.com/gcash/bchd/wire"
	"github.com/gcash/bchutil"
	"github.com/wemeetagain/go-hdwallet"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	exitVal := m.Run()
	os.Exit(exitVal)
}

// This example demonstrates creating a script which pays to a bitcoin address.
// It also prints the created script hex and uses the DisasmString function to
// display the disassembled script.
func Test_Example_(t *testing.T) {
	t.Run("example pay to script", func(t *testing.T) {
		// Parse the address to send the coins to into a bchutil.Address
		// which is useful to ensure the accuracy of the address and determine
		// the address type.  It is also required for the upcoming call to
		// PayToAddrScript.
		addressStr := "bitcoincash:qqfgqp8l9l90zwetj84k2jcac2m8falvvydrpuu45u"
		address, err := bchutil.DecodeAddress(addressStr, &chaincfg.MainNetParams)
		assert.NoError(t, err)

		// Create a public key script that pays to the address.
		script, err := txscript.PayToAddrScript(address)
		assert.NoError(t, err)
		fmt.Printf("Script Hex: %x\n", script)

		disasm, err := txscript.DisasmString(script)
		assert.NoError(t, err)
		fmt.Println("Script Disassembly:", disasm)

		// Output:
		// Script Hex: 76a914128004ff2fcaf13b2b91eb654b1dc2b674f7ec6188ac
		// Script Disassembly: OP_DUP OP_HASH160 128004ff2fcaf13b2b91eb654b1dc2b674f7ec61 OP_EQUALVERIFY OP_CHECKSIG
	})

	t.Run("example extract pk script addrs", func(t *testing.T) {
		// Start with a standard pay-to-pubkey-hash script.
		scriptHex := "76a914128004ff2fcaf13b2b91eb654b1dc2b674f7ec6188ac"
		script, err := hex.DecodeString(scriptHex)
		assert.NoError(t, err)

		// Extract and print details from the script.
		scriptClass, addresses, reqSigs, err := txscript.ExtractPkScriptAddrs(
			script, &chaincfg.MainNetParams)
		assert.NoError(t, err)
		fmt.Println("Script Class:", scriptClass)
		fmt.Println("Addresses:", addresses)
		fmt.Println("Required Signatures:", reqSigs)

		// Output:
		// Script Class: pubkeyhash
		// Addresses: [qqfgqp8l9l90zwetj84k2jcac2m8falvvydrpuu45u]
		// Required Signatures: 1
	})

	// This example demonstrates manually creating and signing a redeem transaction.
	t.Run("example Sign tx output", func(t *testing.T) {
		// Ordinarily the private key would come from whatever storage mechanism
		// is being used, but for this example just hard code it.
		privKeyBytes, err := hex.DecodeString("22a47fa09a223f2aa079edf85a7c2" +
			"d4f8720ee63e502ee2869afab7de234b80c")
		if err != nil {
			fmt.Println(err)
			return
		}
		privKey, pubKey := bchec.PrivKeyFromBytes(bchec.S256(), privKeyBytes)
		pubKeyHash := bchutil.Hash160(pubKey.SerializeCompressed())
		addr, err := bchutil.NewAddressPubKeyHash(pubKeyHash,
			&chaincfg.MainNetParams)
		assert.NoError(t, err)

		// For this example, create a fake transaction that represents what
		// would ordinarily be the real transaction that is being spent.  It
		// contains a single output that pays to address in the amount of 1 BCH.
		originTx := wire.NewMsgTx(wire.TxVersion)
		prevOut := wire.NewOutPoint(&chainhash.Hash{}, ^uint32(0))
		txIn := wire.NewTxIn(prevOut, []byte{txscript.OP_0, txscript.OP_0})
		originTx.AddTxIn(txIn)
		pkScript, err := txscript.PayToAddrScript(addr)
		assert.NoError(t, err)

		txOut := wire.NewTxOut(100000000, pkScript)
		originTx.AddTxOut(txOut)
		originTxHash := originTx.TxHash()

		// Create the transaction to redeem the fake transaction.
		redeemTx := wire.NewMsgTx(wire.TxVersion)

		// Add the input(s) the redeeming transaction will spend.  There is no
		// signature script at this point since it hasn't been created or signed
		// yet, hence nil is provided for it.
		prevOut = wire.NewOutPoint(&originTxHash, 0)
		txIn = wire.NewTxIn(prevOut, nil)
		redeemTx.AddTxIn(txIn)

		// Ordinarily this would contain that actual destination of the funds,
		// but for this example don't bother.
		txOut = wire.NewTxOut(0, nil)
		redeemTx.AddTxOut(txOut)

		// Sign the redeeming transaction.
		lookupKey := func(a bchutil.Address) (*bchec.PrivateKey, bool, error) {
			// Ordinarily this function would involve looking up the private
			// key for the provided address, but since the only thing being
			// signed in this example uses the address associated with the
			// private key from above, simply return it with the compressed
			// flag set since the address is using the associated compressed
			// public key.
			//
			// NOTE: If you want to prove the code is actually signing the
			// transaction properly, uncomment the following line which
			// intentionally returns an invalid key to sign with, which in
			// turn will result in a failure during the script execution
			// when verifying the signature.
			//
			// privKey.D.SetInt64(12345)
			//
			return privKey, true, nil
		}
		// Notice that the script database parameter is nil here since it isn't
		// used.  It must be specified when pay-to-script-hash transactions are
		// being signed.
		sigScript, err := txscript.SignTxOutput(&chaincfg.MainNetParams,
			redeemTx, 0, -1, originTx.TxOut[0].PkScript, txscript.SigHashAll,
			txscript.KeyClosure(lookupKey), nil, nil)
		assert.NoError(t, err)

		redeemTx.TxIn[0].SignatureScript = sigScript

		// Prove that the transaction has been validly signed by executing the
		// script pair.
		flags := txscript.ScriptBip16 | txscript.ScriptVerifyDERSignatures |
			txscript.ScriptStrictMultiSig |
			txscript.ScriptDiscourageUpgradableNops |
			txscript.ScriptVerifyBip143SigHash |
			txscript.ScriptVerifySchnorr

		vm, err := txscript.NewEngine(originTx.TxOut[0].PkScript, redeemTx, 0,
			flags, nil, nil, -1)
		assert.NoError(t, err)

		err = vm.Execute()
		assert.NoError(t, err)
		fmt.Println("Transaction successfully signed")

		// Output:
		// Transaction successfully signed
	})

	var (
		sign    *Signature
		pubkey  *bchec.PublicKey
		privKey *bchec.PrivateKey
		hash    []byte
	)

	t.Run("example SignSchnorr", func(t *testing.T) {
		seed := crypt.CreateHash()

		mnemonic, err := mnemonic.New([]byte(seed), mnemonic.English)
		assert.NoError(t, err)

		// Create a master private key
		masterPrv := hdwallet.MasterKey([]byte(mnemonic.Sentence()))
		privKey, pubkey = bchec.PrivKeyFromBytes(bchec.S256(), masterPrv.Key)

		message := "Satoshi Nakamoto"
		h := sha256.Sum256([]byte(message))
		hash := h[:]

		sign, err = signSchnorr(privKey, hash)
		assert.NoError(t, err)
		fmt.Println(fmt.Sprintf("**************************** SignSchnorr %+v", sign))
		t.Fail()
	})

	t.Run("example SignSchnorr verify", func(t *testing.T) {
		sig, err := bchec.ParseSchnorrSignature(sign.Serialize())
		if err != nil {
			t.Fatal(err)
		}

		valid := sig.Verify(hash, pubkey)
		assert.True(t, valid)
	})
}

// SignatureType enumerates the type of signature. Either ECDSA or Schnorr
type SignatureType uint8

// Signature is a type representing either an ecdsa or schnorr signature.
type Signature struct {
	R       *big.Int
	S       *big.Int
	sigType SignatureType
}

const (
	// SignatureTypeECDSA defines an ecdsa signature
	SignatureTypeECDSA SignatureType = iota

	// SignatureTypeSchnorr defines a schnorr signature
	SignatureTypeSchnorr
)

var (
	// Used in RFC6979 implementation when testing the nonce for correctness
	one = big.NewInt(1)

	// oneInitializer is used to fill a byte slice with byte 0x01.  It is provided
	// here to avoid the need to create it multiple times.
	oneInitializer = []byte{0x01}
)

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
		R:       rx,
		S:       s,
		sigType: SignatureTypeSchnorr,
	}, nil
}

// nonceRFC6979 generates an ECDSA nonce (`k`) deterministically according to RFC 6979.
// It takes a 32-byte hash as an input and returns 32-byte nonce to be used in ECDSA algorithm.
func nonceRFC6979(privkey *big.Int, hash []byte, additionalData []byte) *big.Int {
	curve := bchec.S256()
	q := curve.Params().N
	x := privkey
	alg := sha256.New

	qlen := q.BitLen()
	holen := alg().Size()
	rolen := (qlen + 7) >> 3
	bx := append(int2octets(x, rolen), bits2octets(hash, curve, rolen)...)

	// Step B
	v := bytes.Repeat(oneInitializer, holen)

	// Step C (Go zeroes the all allocated memory)
	k := make([]byte, holen)

	// Step D
	if additionalData != nil {
		k = mac(alg, k, append(append(append(v, 0x00), bx...), additionalData...))
	} else {
		k = mac(alg, k, append(append(v, 0x00), bx...))
	}

	// Step E
	v = mac(alg, k, v)

	// Step F
	if additionalData != nil {
		k = mac(alg, k, append(append(append(v, 0x01), bx...), additionalData...))
	} else {
		k = mac(alg, k, append(append(v, 0x01), bx...))
	}

	// Step G
	v = mac(alg, k, v)

	// Step H
	for {
		// Step H1
		var t []byte

		// Step H2
		for len(t)*8 < qlen {
			v = mac(alg, k, v)
			t = append(t, v...)
		}

		// Step H3
		secret := hashToInt(t, curve)
		if secret.Cmp(one) >= 0 && secret.Cmp(q) < 0 {
			return secret
		}
		k = mac(alg, k, append(v, 0x00))
		v = mac(alg, k, v)
	}
}

// padIntBytes pads a big int bytes with leading zeros if they
// are missing to get the length up to 32 bytes.
func padIntBytes(val *big.Int) []byte {
	b := val.Bytes()
	pad := bytes.Repeat([]byte{0x00}, 32-len(b))
	return append(pad, b...)
}

// https://tools.ietf.org/html/rfc6979#section-2.3.3
func int2octets(v *big.Int, rolen int) []byte {
	out := v.Bytes()

	// left pad with zeros if it's too short
	if len(out) < rolen {
		out2 := make([]byte, rolen)
		copy(out2[rolen-len(out):], out)
		return out2
	}

	// drop most significant bytes if it's too long
	if len(out) > rolen {
		out2 := make([]byte, rolen)
		copy(out2, out[len(out)-rolen:])
		return out2
	}

	return out
}

// hashToInt converts a hash value to an integer. There is some disagreement
// about how this is done. [NSA] suggests that this is done in the obvious
// manner, but [SECG] truncates the hash to the bit-length of the curve order
// first. We follow [SECG] because that's what OpenSSL does. Additionally,
// OpenSSL right shifts excess bits from the number if the hash is too large
// and we mirror that too.
// This is borrowed from crypto/ecdsa.
func hashToInt(hash []byte, c elliptic.Curve) *big.Int {
	orderBits := c.Params().N.BitLen()
	orderBytes := (orderBits + 7) / 8
	if len(hash) > orderBytes {
		hash = hash[:orderBytes]
	}

	ret := new(big.Int).SetBytes(hash)
	excess := len(hash)*8 - orderBits
	if excess > 0 {
		ret.Rsh(ret, uint(excess))
	}
	return ret
}

// mac returns an HMAC of the given key and message.
func mac(alg func() hash.Hash, k, m []byte) []byte {
	h := hmac.New(alg, k)
	h.Write(m)
	return h.Sum(nil)
}

// https://tools.ietf.org/html/rfc6979#section-2.3.4
func bits2octets(in []byte, curve elliptic.Curve, rolen int) []byte {
	z1 := hashToInt(in, curve)
	z2 := new(big.Int).Sub(z1, curve.Params().N)
	if z2.Sign() < 0 {
		return int2octets(z1, rolen)
	}
	return int2octets(z2, rolen)
}

// Serialize returns the a serialized signature depending on the SignatureType.
// Note that the serialized bytes returned do not include the appended hash type
// used in Bitcoin signature scripts.
//
// ECDSA signature in the more strict DER format.
//
// encoding/asn1 is broken so we hand roll this output:
//
// 0x30 <length> 0x02 <length r> r 0x02 <length s> s
func (sig *Signature) Serialize() []byte {
	// Schnorr signatures are easy to serialize
	if sig.sigType == SignatureTypeSchnorr {
		return append(padIntBytes(sig.R), padIntBytes(sig.S)...)
	}
	// low 'S' malleability breaker
	sigS := sig.S

	halfOrder := new(big.Int).Rsh(bchec.S256().N, 1)
	if sigS.Cmp(halfOrder) == 1 {
		sigS = new(big.Int).Sub(bchec.S256().N, sigS)
	}
	// Ensure the encoded bytes for the r and s values are canonical and
	// thus suitable for DER encoding.
	rb := canonicalizeInt(sig.R)
	sb := canonicalizeInt(sigS)

	// total length of returned signature is 1 byte for each magic and
	// length (6 total), plus lengths of r and s
	length := 6 + len(rb) + len(sb)
	b := make([]byte, length)

	b[0] = 0x30
	b[1] = byte(length - 2)
	b[2] = 0x02
	b[3] = byte(len(rb))
	offset := copy(b[4:], rb) + 4
	b[offset] = 0x02
	b[offset+1] = byte(len(sb))
	copy(b[offset+2:], sb)
	return b
}

// canonicalizeInt returns the bytes for the passed big integer adjusted as
// necessary to ensure that a big-endian encoded integer can't possibly be
// misinterpreted as a negative number.  This can happen when the most
// significant bit is set, so it is padded by a leading zero byte in this case.
// Also, the returned bytes will have at least a single byte when the passed
// value is 0.  This is required for DER encoding.
func canonicalizeInt(val *big.Int) []byte {
	b := val.Bytes()
	if len(b) == 0 {
		b = []byte{0x00}
	}
	if b[0]&0x80 != 0 {
		paddedBytes := make([]byte, len(b)+1)
		copy(paddedBytes[1:], b)
		b = paddedBytes
	}
	return b
}
