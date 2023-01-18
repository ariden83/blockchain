package blockchain

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Genesis(t *testing.T) {
	transaction := &Transaction{ID: []byte("id-1")}
	firstBlock := Genesis(transaction)
	assert.NotNil(t, firstBlock)
}

func Test_calculateHash(t *testing.T) {
	hash := calculateHash(Block{Index: big.NewInt(3)})
	assert.Equal(t, []byte{0x36, 0x32, 0x34, 0x62, 0x36, 0x30, 0x63, 0x35, 0x38, 0x63, 0x39, 0x64, 0x38, 0x62, 0x66, 0x62, 0x36, 0x66, 0x66, 0x31, 0x38, 0x38, 0x36, 0x63, 0x32, 0x66, 0x64, 0x36, 0x30, 0x35, 0x64, 0x32, 0x61, 0x64, 0x65, 0x62, 0x36, 0x65, 0x61, 0x34, 0x64, 0x61, 0x35, 0x37, 0x36, 0x30, 0x36, 0x38, 0x32, 0x30, 0x31, 0x62, 0x36, 0x63, 0x36, 0x39, 0x35, 0x38, 0x63, 0x65, 0x39, 0x33, 0x66, 0x34}, hash)

	hash = calculateHash(Block{Index: big.NewInt(1)})
	assert.Equal(t, []byte{0x34, 0x61, 0x34, 0x34, 0x64, 0x63, 0x31, 0x35, 0x33, 0x36, 0x34, 0x32, 0x30, 0x34, 0x61, 0x38, 0x30, 0x66, 0x65, 0x38, 0x30, 0x65, 0x39, 0x30, 0x33, 0x39, 0x34, 0x35, 0x35, 0x63, 0x63, 0x31, 0x36, 0x30, 0x38, 0x32, 0x38, 0x31, 0x38, 0x32, 0x30, 0x66, 0x65, 0x32, 0x62, 0x32, 0x34, 0x66, 0x31, 0x65, 0x35, 0x32, 0x33, 0x33, 0x61, 0x64, 0x65, 0x36, 0x61, 0x66, 0x31, 0x64, 0x64, 0x35}, hash)
}

func Test_NextID(t *testing.T) {
	for _, test := range []struct {
		providedValue *big.Int
		expected      *big.Int
	}{
		{providedValue: big.NewInt(0), expected: big.NewInt(1)},
		{providedValue: big.NewInt(1), expected: big.NewInt(2)},
		{providedValue: big.NewInt(2), expected: big.NewInt(3)},
		{providedValue: big.NewInt(10), expected: big.NewInt(11)},
	} {
		assert.Equal(t, test.expected, NextID(test.providedValue))
	}
}

func Test_GetLastBlock(t *testing.T) {
	lastBlock := Block{Index: big.NewInt(3)}

	BlockChain = []Block{{
		Index: big.NewInt(1),
	}, {
		Index: big.NewInt(2),
	}, lastBlock}

	block := GetLastBlock()
	assert.Equal(t, lastBlock, block)
}

func Test_isHashValid(t *testing.T) {
	for _, test := range []struct {
		providedDifficulty int
		providedValue      string
		expected           bool
	}{
		{providedDifficulty: 1, providedValue: "0", expected: true},
		{providedDifficulty: 1, providedValue: "9876543", expected: false},
		{providedDifficulty: 1, providedValue: "09876543", expected: true},
		{providedDifficulty: 2, providedValue: "9876543", expected: false},
		{providedDifficulty: 2, providedValue: "09876543", expected: false},
		{providedDifficulty: 2, providedValue: "009876543", expected: true},
		{providedDifficulty: 3, providedValue: "9876543", expected: false},
		{providedDifficulty: 3, providedValue: "09876543", expected: false},
		{providedDifficulty: 3, providedValue: "009876543", expected: false},
		{providedDifficulty: 3, providedValue: "0009876543", expected: true},
	} {
		assert.Equal(t, test.expected, isHashValid([]byte(test.providedValue), test.providedDifficulty))
	}
}

func Test_IsValid(t *testing.T) {
	listBlock := []Block{{
		Index: big.NewInt(1),
	}, {
		Index: big.NewInt(2),
	}}

	assert.True(t, IsValid(listBlock))
}
