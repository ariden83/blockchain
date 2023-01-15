package utils

import (
	"math/big"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ariden83/blockchain/internal/blockchain"
)

func Test_Serialize_Deserialize(t *testing.T) {
	toSerialize := []blockchain.Block{
		{
			Index: big.NewInt(1),
		},
	}
	serializedStr, err := Serialize(toSerialize)
	assert.NoError(t, err)
	assert.NotEmpty(t, serializedStr)

	block := &[]blockchain.Block{}
	err = Deserialize(serializedStr, block)
	assert.NoError(t, err)
	assert.Len(t, *block, 1)
}

func Test_RandomString(t *testing.T) {
	randomString := RandomString(5)
	list := strings.Split(randomString, " ")
	assert.NotNil(t, len(list), 5)

	randomString = RandomString(0)
	list = strings.Split(randomString, " ")
	assert.NotNil(t, len(list), 0)
}
