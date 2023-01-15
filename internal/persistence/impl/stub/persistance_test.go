package stub

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Persistence(t *testing.T) {
	p := New()
	t.Run("New", func(t *testing.T) {
		assert.NotNil(t, p)
	})

	t.Run("GetLastHash", func(t *testing.T) {
		lastHash, err := p.GetLastHash()
		assert.NoError(t, err)
		assert.Nil(t, lastHash)
	})

	t.Run("Update", func(t *testing.T) {
		lastHash := []byte("hash")
		err := p.Update(lastHash, lastHash)
		assert.NoError(t, err)
	})

	t.Run("LastHash", func(t *testing.T) {
		lastHash := p.LastHash()
		assert.Nil(t, lastHash)
	})

	t.Run("GetCurrentHashSerialize", func(t *testing.T) {
		hashSerialize := []byte("hashSerialize")
		lastHash, err := p.GetCurrentHashSerialize(hashSerialize)
		assert.NoError(t, err)
		assert.Nil(t, lastHash)
	})

	t.Run("GetCurrentHashSerialize", func(t *testing.T) {
		assert.True(t, p.DBExists())
	})

	t.Run("SetLastHash", func(t *testing.T) {
		lastHash := []byte("last hash set")
		p.SetLastHash(lastHash)
		assert.Equal(t, p.lastHash, lastHash)
	})

	t.Run("Close", func(t *testing.T) {
		assert.Nil(t, p.Close())
	})
}
