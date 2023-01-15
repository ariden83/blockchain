package validator

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ariden83/blockchain/internal/blockchain"
)

func Test_Validator(t *testing.T) {
	validater := New(blockchain.Block{}, []string{})

	t.Run("New", func(t *testing.T) {
		assert.NotNil(t, validater)
	})

	t.Run("IsAcceptedByMajority", func(t *testing.T) {
		assert.False(t, validater.IsAcceptedByMajority())
	})

	t.Run("IsRefusedByMajority", func(t *testing.T) {
		assert.False(t, validater.IsRefusedByMajority())
	})

	t.Run("Accept", func(t *testing.T) {
		validater.Accept()
		assert.True(t, validater.IsAcceptedByMajority())
	})

	t.Run("Refuse", func(t *testing.T) {
		validater.Refuse()
		assert.True(t, validater.IsRefusedByMajority())
	})
}
