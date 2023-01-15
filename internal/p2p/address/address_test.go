package address

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Address(t *testing.T) {
	add := New()
	me := "my-address"

	t.Run("Create a new address", func(t *testing.T) {
		assert.NotNil(t, add)
	})

	t.Run("Me must be empty", func(t *testing.T) {
		assert.Empty(t, add.Address())
	})

	t.Run("SetIAM", func(t *testing.T) {
		add.SetAddress(me)
		assert.Equal(t, me, add.Address())
	})

	t.Run("RecreateAddress", func(t *testing.T) {
		addressRecreated := add.RecreateMyAddress()

		assert.Equal(t, "[\"my-address\"]", string(addressRecreated))
	})

	t.Run("CurrentAddress", func(t *testing.T) {
		currentAddress := add.CurrentAddress()
		assert.Equal(t, []string{me}, currentAddress)
	})
}
