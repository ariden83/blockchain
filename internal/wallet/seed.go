package wallet

import (
	"github.com/LuisAcerv/btchdwallet/crypt"
	"github.com/brianium/mnemonic"
	"github.com/wemeetagain/go-hdwallet"
	"sync"
	"time"
)

// Seed represents each 'item' in the blockchain
type Seed struct {
	Address   string
	Timestamp string
	PubKey    string
	PrivKey   string
	Mnemonic  string
}

type SeedNoPrivKey struct {
	Address   string
	Timestamp string
	PubKey    string
}

var mutex = &sync.Mutex{}

func (ws *Wallets) GetAllSeeds() []SeedNoPrivKey {
	var allSeeds []SeedNoPrivKey
	for _, j := range ws.Seeds {
		allSeeds = append(allSeeds, SeedNoPrivKey{
			Address:   j.Address,
			Timestamp: j.Timestamp,
			PubKey:    j.PubKey,
		})
	}
	return allSeeds
}

func (ws *Wallets) Create() (*Seed, error) {

	seed := crypt.CreateHash()
	mnemonic, err := mnemonic.New([]byte(seed), mnemonic.English)
	if err != nil {
		return nil, err
	}

	// Create a master private key
	masterPrv := hdwallet.MasterKey([]byte(mnemonic.Sentence()))

	// Convert a private key to public key
	masterPub := masterPrv.Pub()

	// Get your address
	address := masterPub.Address()

	t := time.Now()
	newSeed := Seed{
		Address:   address,
		PubKey:    masterPub.String(),
		PrivKey:   masterPrv.String(),
		Mnemonic:  mnemonic.Sentence(),
		Timestamp: t.String(),
	}

	mutex.Lock()
	ws.Seeds = append(ws.Seeds, newSeed)
	mutex.Unlock()

	ws.Save()

	return &newSeed, nil
}
