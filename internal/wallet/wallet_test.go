package wallet

import (
	"log"
	"testing"

	"github.com/LuisAcerv/btchdwallet/crypt"
	"github.com/brianium/mnemonic"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/ariden83/blockchain/config"
	"github.com/ariden83/blockchain/internal/hdwallet"
	"github.com/ariden83/blockchain/internal/logger"
)

func Test_New_Wallet(t *testing.T) {
	wallerAdapter, err := New(Config{}, zap.NewNop())
	assert.NoError(t, err)
	assert.NotNil(t, wallerAdapter)
}

func Test_Hash(t *testing.T) {
	t.Run("test hash comparaison", func(t *testing.T) {
		mnemonic := []byte("couple robot escape silent main once smoke check good basket mimic similar")
		mnemonicHash := hash(mnemonic)
		for i := 1; i < 10; i++ {
			mnemonicCurrent := hash(mnemonic)
			// require.NoError(t, err)
			assert.Equal(t, mnemonicCurrent, mnemonicHash)
		}
	})
}

func Test_HDWallet_PubKey_Always_Same(t *testing.T) {
	seed := crypt.CreateHash()
	mnemonic, err := mnemonic.New([]byte(seed), mnemonic.English)
	assert.NoError(t, err)

	masterPrv := hdwallet.MasterKey([]byte(mnemonic.Sentence()))
	masterPub := masterPrv.Pub()

	for i := 1; i < 100; i++ {
		masterPrvTest := hdwallet.MasterKey([]byte(mnemonic.Sentence()))
		assert.Equal(t, masterPrv.String(), masterPrvTest.String())
		assert.Equal(t, masterPub.String(), masterPrvTest.Pub().String())
	}
}

func Test_HDWallet_PrivKey_String_PrivKey(t *testing.T) {
	seed := crypt.CreateHash()
	mnemonic, err := mnemonic.New([]byte(seed), mnemonic.English)
	assert.NoError(t, err)

	masterPrv := hdwallet.MasterKey([]byte(mnemonic.Sentence()))
	masterStrPub := masterPrv.Pub().String()
	strPrivKey := masterPrv.String()

	for i := 1; i < 100; i++ {
		masterPrvTest, err := hdwallet.StringWallet(strPrivKey)
		assert.NoError(t, err)
		assert.Equal(t, masterStrPub, masterPrvTest.Pub().String())
	}
}

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
	password := []byte("123456")
	seedCreate, err := w.Create(password)
	require.NoError(t, err)
	require.NotNil(t, seedCreate)
	require.NotEmpty(t, seedCreate.Mnemonic)
	require.NotEmpty(t, seedCreate.PubKey)
	require.NotEmpty(t, seedCreate.Address)

	isValidate := w.Validate(seedCreate.PrivKey)
	require.True(t, isValidate)

	seed, err := w.Seed(seedCreate.Mnemonic, password)
	require.NoError(t, err)
	require.NotNil(t, seed)
	require.Equal(t, seed.PubKey, seedCreate.PubKey)
	require.Equal(t, seed.Address, seedCreate.Address)
}
