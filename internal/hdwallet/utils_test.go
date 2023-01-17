package hdwallet

import (
	"math/big"
	"testing"

	"github.com/btcsuite/btcutil/base58"
	"github.com/stretchr/testify/assert"
)

func Test_S256(t *testing.T) {
	s256 := S256()
	assert.NotNil(t, s256)
	if !s256.IsOnCurve(s256.Params().Gx, s256.Params().Gy) {
		t.Fatal("generator point does not claim to be on the curve")
	}
}

func Test_hash160(t *testing.T) {
	t.Run("must be ok with nil value", func(t *testing.T) {
		pubKey := hash160(nil)
		assert.Equal(t, []byte{0xb4, 0x72, 0xa2, 0x66, 0xd0, 0xbd, 0x89, 0xc1, 0x37, 0x6, 0xa4, 0x13, 0x2c, 0xcf, 0xb1, 0x6f, 0x7c, 0x3b, 0x9f, 0xcb}, pubKey)

	})

	t.Run("must be ok with empty value", func(t *testing.T) {
		pubKey := hash160([]byte(""))
		assert.Equal(t, []byte{0xb4, 0x72, 0xa2, 0x66, 0xd0, 0xbd, 0x89, 0xc1, 0x37, 0x6, 0xa4, 0x13, 0x2c, 0xcf, 0xb1, 0x6f, 0x7c, 0x3b, 0x9f, 0xcb}, pubKey)

	})

	t.Run("must be ok", func(t *testing.T) {
		hash160 := hash160([]byte("data-to-hash"))
		assert.Equal(t, []byte{0x6, 0xe4, 0x1c, 0xc, 0xa2, 0xff, 0x4b, 0x8, 0x92, 0x3e, 0x65, 0x28, 0x38, 0xb7, 0xb, 0x16, 0xa3, 0xe1, 0x89, 0xab}, hash160)
	})
}

func Test_dblSha256(t *testing.T) {
	t.Run("must be ok with nil value", func(t *testing.T) {
		pubKey := dblSha256(nil)
		assert.Equal(t, []byte{0x5d, 0xf6, 0xe0, 0xe2, 0x76, 0x13, 0x59, 0xd3, 0xa, 0x82, 0x75, 0x5, 0x8e, 0x29, 0x9f, 0xcc, 0x3, 0x81, 0x53, 0x45, 0x45, 0xf5, 0x5c, 0xf4, 0x3e, 0x41, 0x98, 0x3f, 0x5d, 0x4c, 0x94, 0x56}, pubKey)

	})

	t.Run("must be ok with empty value", func(t *testing.T) {
		pubKey := dblSha256([]byte(""))
		assert.Equal(t, []byte{0x5d, 0xf6, 0xe0, 0xe2, 0x76, 0x13, 0x59, 0xd3, 0xa, 0x82, 0x75, 0x5, 0x8e, 0x29, 0x9f, 0xcc, 0x3, 0x81, 0x53, 0x45, 0x45, 0xf5, 0x5c, 0xf4, 0x3e, 0x41, 0x98, 0x3f, 0x5d, 0x4c, 0x94, 0x56}, pubKey)

	})

	t.Run("must be ok", func(t *testing.T) {
		dblSha256 := dblSha256([]byte("data-to-hash"))
		assert.Equal(t, []byte{0xd3, 0xcc, 0x1a, 0x21, 0xf5, 0xbd, 0x14, 0xeb, 0xa7, 0x24, 0xf0, 0x55, 0xb0, 0x70, 0xa7, 0xf1, 0x4c, 0x13, 0x7e, 0x77, 0x6f, 0xad, 0x18, 0xd4, 0x5f, 0xd3, 0x8b, 0x67, 0xc3, 0x92, 0xf0, 0xf6}, dblSha256)
	})
}

func Test_privToPub(t *testing.T) {
	t.Run("must be ok with nil value", func(t *testing.T) {
		pubKey := privToPub(nil)
		assert.Equal(t, []byte{0x2, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}, pubKey)

	})

	t.Run("must be ok with empty value", func(t *testing.T) {
		pubKey := privToPub([]byte(""))
		assert.Equal(t, []byte{0x2, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}, pubKey)

	})

	t.Run("must be ok", func(t *testing.T) {
		pubKey := privToPub([]byte("data-to-hash"))
		assert.Equal(t, []byte{0x3, 0xcc, 0xb8, 0x3, 0x82, 0xb9, 0x73, 0x46, 0x1e, 0x44, 0xb, 0x38, 0x7a, 0xa8, 0x24, 0x16, 0x22, 0x14, 0xdc, 0x5a, 0xe6, 0x4, 0x21, 0x86, 0xa4, 0x6, 0x1c, 0xf4, 0x43, 0xe8, 0x10, 0x6a, 0xfd}, pubKey)
	})
}

func Test_onCurve(t *testing.T) {
	t.Run("must be false with x nil value", func(t *testing.T) {
		y := new(big.Int)

		isOnCurve := onCurve(nil, y)
		assert.False(t, isOnCurve)
	})

	t.Run("must be false with y nil value", func(t *testing.T) {
		x := new(big.Int)

		isOnCurve := onCurve(x, nil)
		assert.False(t, isOnCurve)
	})

	t.Run("must be ok", func(t *testing.T) {
		data := base58.Decode(base58.Encode([]byte("00eb15231dfceb60925886b67d065299925915aeb172c06647")))
		x, y := expand(data)
		assert.NotNil(t, x)
		assert.NotNil(t, y)

		isOnCurve := onCurve(x, y)
		assert.True(t, isOnCurve)
	})

	t.Run("must fail with invalid x y value", func(t *testing.T) {
		x, ok := new(big.Int).SetString("218882420012223", 10)
		assert.True(t, ok)

		y, ok := new(big.Int).SetString("2188824200011112223", 10)
		assert.True(t, ok)

		isOnCurve := onCurve(x, y)
		assert.False(t, isOnCurve)
	})
}

func Test_compress(t *testing.T) {

	t.Run("must be false with x nil value", func(t *testing.T) {
		y, ok := new(big.Int).SetString("218882420012223", 10)
		assert.True(t, ok)

		compress := compress(nil, y)
		assert.Empty(t, compress)
	})

	t.Run("must be false with y nil value", func(t *testing.T) {
		x, ok := new(big.Int).SetString("218882420012223", 10)
		assert.True(t, ok)

		compress := compress(x, nil)
		assert.Empty(t, compress)
	})

	t.Run("must be ok", func(t *testing.T) {
		x, ok := new(big.Int).SetString("218882420012223", 10)
		assert.True(t, ok)

		y, ok := new(big.Int).SetString("2188824200011112223", 10)
		assert.True(t, ok)

		compress := compress(x, y)
		assert.NotNil(t, compress)
		assert.Equal(t, []byte{0x3, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xc7, 0x12, 0x88, 0xe4, 0x74, 0xbf}, compress)
	})
}

func Test_expand(t *testing.T) {
	t.Run("must be ok with pub key", func(t *testing.T) {
		pubKey := privToPub([]byte("data-to-hash"))
		x, y := expand(pubKey)

		_, ok := new(big.Int).SetString("75068814825383104477272410521391746915676718919767125217262038807184606888987", 10)
		assert.True(t, ok)

		_, ok = new(big.Int).SetString("75068814825383104477272410521391746915676718919767125217262038807184606888987", 10)
		assert.True(t, ok)

		assert.NotNil(t, x)
		assert.NotNil(t, y)
	})

	t.Run("must return nil values with nil value", func(t *testing.T) {
		x, y := expand(nil)

		assert.Nil(t, x)
		assert.Nil(t, y)
	})

	t.Run("must return nil values with empty byte", func(t *testing.T) {
		x, y := expand([]byte(""))

		assert.Nil(t, x)
		assert.Nil(t, y)
	})

	t.Run("must return nil values with value 1", func(t *testing.T) {
		x, y := expand([]byte("1"))

		assert.Nil(t, x)
		assert.Nil(t, y)
	})

	t.Run("must return nil values with value hash160", func(t *testing.T) {
		data := []byte("data-to-hash")
		hash160 := hash160(data)

		x, y := expand(hash160)

		assert.Nil(t, x)
		assert.Nil(t, y)
	})

	t.Run("must be ok nil values with value dblSha256", func(t *testing.T) {
		data := []byte("data-to-hash")
		dblSha256 := dblSha256(data)

		x, y := expand(dblSha256)

		assert.Nil(t, x)
		assert.Nil(t, y)
	})
}

func Test_addPrivKeys(t *testing.T) {
	t.Run("must be ok with empty bytes", func(t *testing.T) {
		privKey := addPrivKeys([]byte(""), []byte(""))

		assert.NotNil(t, privKey)
		assert.Equal(t, []byte{0x0}, privKey)
	})

	t.Run("must be ok with nil bytes", func(t *testing.T) {
		privKey := addPrivKeys(nil, nil)

		assert.NotNil(t, privKey)
		assert.Equal(t, []byte{0x0}, privKey)
	})
}

/*
func Test_addPubKeys(t *testing.T) {
}

func Test_uint32ToByte(t *testing.T) {
}

func Test_uint16ToByte(t *testing.T) {
}

func Test_byteToUint16(t *testing.T) {
}
*/
