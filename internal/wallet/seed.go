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
	Timestamp int64
	PubKey    string
	PrivKey   string
	Mnemonic  string
}

type SeedNoPrivKey struct {
	Timestamp int64
	PubKey    string
	PrivKey   string
}

var mutex = &sync.Mutex{}

func (ws *Wallets) GetSeeds() *[]Seed {
	return &ws.Seeds
}

func (ws *Wallets) GetAllPublicSeeds() []SeedNoPrivKey {
	var allSeeds []SeedNoPrivKey
	for _, j := range ws.Seeds {
		allSeeds = append(allSeeds, SeedNoPrivKey{
			PrivKey:   j.PrivKey,
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

	t := time.Now().UnixNano() / int64(time.Millisecond)
	newSeed := Seed{
		Address:   address,
		PubKey:    masterPub.String(),
		PrivKey:   masterPrv.String(),
		Mnemonic:  mnemonic.Sentence(),
		Timestamp: t,
	}

	mutex.Lock()
	ws.Seeds = append(ws.Seeds, newSeed)
	mutex.Unlock()

	ws.Save()

	return &newSeed, nil
}
