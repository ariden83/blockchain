package genesis

import (
	"errors"
	"fmt"
	"sync"

	"github.com/davecgh/go-spew/spew"

	"github.com/ariden83/blockchain/config"
	"github.com/ariden83/blockchain/internal/blockchain"
	"github.com/ariden83/blockchain/internal/event"
	"github.com/ariden83/blockchain/internal/p2p"
	"github.com/ariden83/blockchain/internal/persistence"
	"github.com/ariden83/blockchain/internal/transaction"
	"github.com/ariden83/blockchain/internal/utils"
	"github.com/ariden83/blockchain/internal/wallet"
)

var mutex = &sync.Mutex{}

type Genesis struct {
	persistence persistenceadapter.Adapter
	transaction transaction.Adapter
	cfg         *config.Config
	p2p         *p2p.EndPoint
	event       *event.Event
	wallets     wallet.IWallets
}

func New(cfg *config.Config, persistence persistenceadapter.Adapter, trans transaction.Adapter, p *p2p.EndPoint,
	evt *event.Event, wallets wallet.IWallets) *Genesis {
	return &Genesis{
		wallets:     wallets,
		persistence: persistence,
		transaction: trans,
		cfg:         cfg,
		p2p:         p,
		event:       evt,
	}
}

func (g *Genesis) Genesis() bool {
	if g.p2p.Target() == "" {
		// call default genesis
		return false
	}
	return true
}

func (g *Genesis) Load(stop chan error) {
	// si y'a une instance, on la load
	if g.p2p.Enabled() && g.p2p.HasTarget() {
		// on notifie la demande de récupération des fichiers
		g.p2p.PushMsgForFiles(stop)
		return
	}

	// sinon, on créé le premier hash
	var lastHash []byte

	// si les fichiers locaux n'existent pas
	if !g.persistence.DBExists() {
		stop <- errors.New("fail to open DB files")
		return
	}

	lastHash, err := g.persistence.GetLastHash()
	if err != nil {
		stop <- errors.New("fail to get last hash")
		return
	}

	if lastHash == nil {
		lastHash = g.createGenesis(stop)

	} else {

		val, err := g.persistence.GetCurrentHashSerialize(lastHash)
		if err != nil {
			stop <- errors.New("fail to get current hash")
			return
		}

		block := &blockchain.Block{}
		if err := utils.Deserialize(val, block); err != nil {
			stop <- fmt.Errorf("fail to deserialize hash serializes: %w", err)
			return
		}

		g.persistence.SetLastHash(lastHash)

		mutex.Lock()
		blockchain.BlockChain = append(blockchain.BlockChain, *block)
		mutex.Unlock()

		spew.Dump(blockchain.BlockChain)
	}
}

func (g *Genesis) createGenesis(stop chan error) []byte {
	privateKey := []byte(g.cfg.Transactions.PrivateKey)
	if g.cfg.Transactions.PrivateKey == "" {
		seed, _ := g.wallets.Create([]byte("test"))
		privateKey = seed.PrivKey
	}
	g.wallets = nil

	cbtx := g.transaction.CoinBaseTx(privateKey)
	genesis := blockchain.Genesis(cbtx)
	fmt.Println("Genesis proved")

	firstHash := genesis.Hash

	serializeBLock, err := utils.Serialize(genesis)
	if err != nil {
		stop <- fmt.Errorf("fail to serialize genesis: %w", err)
		return []byte{}
	}

	err = g.persistence.Update(firstHash, serializeBLock)
	if err != nil {
		stop <- errors.New("fail to update persistence")
		return []byte{}
	}
	return firstHash
}
