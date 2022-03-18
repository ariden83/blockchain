package signschnorr

import (
	"bytes"
	"crypto/elliptic"
	"crypto/hmac"
	"crypto/sha256"
	"hash"
	"math/big"

	"github.com/gcash/bchd/bchec"
)

// Signature is a type representing either an ecdsa or schnorr signature.
type Signature struct {
	R *big.Int
	S *big.Int
}

var (
	// Used in RFC6979 implementation when testing the nonce for correctness
	one = big.NewInt(1)

	// oneInitializer is used to fill a byte slice with byte 0x01.  It is provided
	// here to avoid the need to create it multiple times.
	oneInitializer = []byte{0x01}
)

// signSchnorr signs the hash using the schnorr signature algorithm.
// The rfc6979 nonce derivation function accepts additional entropy.
// We are using the same entropy that is used by bitcoin-abc so our test
// vectors will be compatible. This byte string is chosen to avoid collisions
// with ECDSA which would render the signature insecure.
//
// See https://github.com/bitcoincashorg/bitcoincash.org/blob/master/spec/2019-05-15-schnorr.md#recommended-practices-for-secure-signature-generation
func SignSchnorr(privateKey *bchec.PrivateKey, hash []byte) (*Signature, error) {
	additionalData := []byte{'S', 'c', 'h', 'n', 'o', 'r', 'r', '+', 'S', 'H', 'A', '2', '5', '6', ' ', '+'}
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

func ParseSchnorrSig(sign *Signature, hash []byte, pubkey *bchec.PublicKey) bool {
	sig, err := bchec.ParseSchnorrSignature(sign.Serialize())
	if err != nil {
		return false
	}
	return sig.Verify(hash, pubkey)
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
	return append(padIntBytes(sig.R), padIntBytes(sig.S)...)
}
