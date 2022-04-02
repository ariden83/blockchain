package wallet

import (
	"go.uber.org/zap"
	"time"

	"github.com/wemeetagain/go-hdwallet"
)

func (w *Wallets) allKeysFromMnemonic(mnemonic []byte) *Seed {
	// Create a master private key
	masterPrv := hdwallet.MasterKey(mnemonic)
	// Convert a private key to public key
	masterPub := masterPrv.Pub()
	// Get your address
	address := masterPub.Address()

	return &Seed{
		PrivKey:   []byte(masterPrv.String()),
		PubKey:    []byte(masterPub.String()),
		Address:   []byte(address),
		Mnemonic:  mnemonic,
		Timestamp: time.Now().UnixNano() / int64(time.Millisecond),
	}
}

func (w *Wallets) GetUserAddress(privKey []byte) string {
	s, err := w.allKeysFromPrivate(privKey)
	if err != nil {
		w.log.Error("fail to get user address", zap.Error(err))
		return ""
	}
	return string(s.Address)
}

func (w *Wallets) allKeysFromPrivate(privKey []byte) (*Seed, error) {
	masterPrv, err := hdwallet.StringWallet(string(privKey))
	if err != nil {
		return nil, err
	}
	// Convert a private key to public key
	masterPub := masterPrv.Pub()
	// Get your address
	address := masterPub.Address()

	return &Seed{
		PrivKey: []byte(masterPrv.String()),
		PubKey:  []byte(masterPub.String()),
		Address: []byte(address),
	}, nil
}
