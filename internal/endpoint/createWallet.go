package endpoint

import (
	"github.com/LuisAcerv/btchdwallet/crypt"
	"github.com/ariden83/blockchain/internal/handle"
	"github.com/ariden83/blockchain/internal/wallet"
	"github.com/brianium/mnemonic"
	"github.com/wemeetagain/go-hdwallet"
	"net/http"
	"time"
)

func (e *EndPoint) handleCreateWallet(w http.ResponseWriter, r *http.Request) {

	seed := crypt.CreateHash()
	mnemonic, err := mnemonic.New([]byte(seed), mnemonic.English)
	handle.Handle(err)

	// Create a master private key
	masterprv := hdwallet.MasterKey([]byte(mnemonic.Sentence()))

	// Convert a private key to public key
	masterpub := masterprv.Pub()

	// Get your address
	address := masterpub.Address()

	t := time.Now()
	newSeed := wallet.Seed{
		Address:   address,
		PubKey:    masterpub.String(),
		PrivKey:   masterprv.String(),
		Mnemonic:  mnemonic.Sentence(),
		Timestamp: t.String(),
	}
	mutex.Lock()
	e.wallets.Seeds = append(e.wallets.Seeds, newSeed)
	mutex.Unlock()

	e.wallets.Save()

	respondWithJSON(w, r, http.StatusCreated, newSeed)
}
