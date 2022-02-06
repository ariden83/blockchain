package wallet

import (
	"log"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ariden83/blockchain/config"
	"github.com/ariden83/blockchain/internal/logger"
)

func Test_Create(t *testing.T) {
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("fail to init persistence %s", err)
	}
	cfg.Wallet.WithFile = false
	cfg.Log.WithFile = false
	logs := logger.InitLight(cfg.Log)
	defer logs.Sync()

	w := Wallets{
		log: logs,
	}
	password := []byte("my-password")
	seedCreate, err := w.Create(password)
	require.NoError(t, err)
	require.NotNil(t, seedCreate)
	require.NotEmpty(t, seedCreate.Mnemonic)
	require.NotEmpty(t, seedCreate.PubKey)
	require.NotEmpty(t, seedCreate.Address)

	seed, err := w.GetSeed([]byte(seedCreate.Mnemonic), password)
	require.NoError(t, err)
	require.NotNil(t, seed)
	require.Equal(t, seed.Mnemonic, seedCreate.Mnemonic)
	require.Equal(t, seed.PubKey, seedCreate.PubKey)
	require.Equal(t, seed.Address, seedCreate.Address)
}
