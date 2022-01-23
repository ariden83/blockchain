package genesis

import (
	"fmt"
	"github.com/ariden83/blockchain/config"
	"github.com/ariden83/blockchain/internal/blockchain"
	"github.com/ariden83/blockchain/internal/event"
	"github.com/ariden83/blockchain/internal/p2p"
	"github.com/ariden83/blockchain/internal/persistence"
	"github.com/ariden83/blockchain/internal/transactions"
	"github.com/ariden83/blockchain/internal/utils"
	"github.com/davecgh/go-spew/spew"
	"sync"
)

var mutex = &sync.Mutex{}

type Genesis struct {
	persistence *persistence.Persistence
	transaction *transactions.Transactions
	cfg         *config.Config
	p2p         *p2p.EndPoint
	event       *event.Event
}

func New(cfg *config.Config, pers *persistence.Persistence, trans *transactions.Transactions,
	p *p2p.EndPoint, evt *event.Event) *Genesis {
	return &Genesis{
		persistence: pers,
		transaction: trans,
		cfg:         cfg,
		p2p:         p,
		event:       evt,
	}
}

func (g *Genesis) Genesis() bool {
	if g.cfg.P2P.Target == "" {
		// call default genesis
		return false
	}
	return true
}

func (g *Genesis) Load(stop chan error) {
	// si y'a une instance, on la load
	if g.p2p.Enabled() && g.p2p.HasTarget() {
		// on notifie la demande de récupération des fichiers
		g.p2p.PushMsgForFiles()
		return
	}

	// sinon, on créé le premier hash

	var lastHash []byte

	// si les fichiers locaux n'existent pas
	if !g.persistence.DBExists() {
		stop <- fmt.Errorf("fail to open DB files")
		return
	}

	lastHash, err := g.persistence.GetLastHash()
	if err != nil {
		stop <- fmt.Errorf("fail to get last hash")
		return
	}

	if lastHash == nil {
		lastHash = g.createGenesis(stop)

	} else {

		val, err := g.persistence.GetCurrentHashSerialize(lastHash)
		if err != nil {
			stop <- fmt.Errorf("fail to get current hash")
			return
		}

		block, err := utils.DeserializeBlock(val)
		if err != nil {
			stop <- fmt.Errorf("fail to deserialize hash serializesd")
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

	var genesisData string = "First Transaction from Genesis" // This is arbitrary public key for our genesis data
	cbtx := g.transaction.CoinBaseTx(g.cfg.Address, genesisData)
	genesis := blockchain.Genesis(cbtx)
	fmt.Println("Genesis proved")

	firstHash := genesis.Hash

	serializeBLock, err := utils.Serialize(genesis)
	if err != nil {
		stop <- fmt.Errorf("fail to serialize genesis")
		return []byte{}
	}

	err = g.persistence.Update(firstHash, serializeBLock)
	if err != nil {
		stop <- fmt.Errorf("fail to update persistence")
		return []byte{}
	}
	return firstHash
}
