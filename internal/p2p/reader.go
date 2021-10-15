package p2p

import (
	"bufio"
	"bytes"
	"encoding/json"
	"github.com/ariden83/blockchain/internal/blockchain"
	"github.com/ariden83/blockchain/internal/event"
	"github.com/ariden83/blockchain/internal/utils"
	"github.com/davecgh/go-spew/spew"
	"go.uber.org/zap"
	"log"
)

// routine Go qui récupère le dernier état de notre blockchain toutes les 5 secondes
// err = rw.Flush()
func (e *EndPoint) readData(rw *bufio.ReadWriter) {

	for {
		str, err := rw.ReadString('\n')
		if err != nil {
			// stream is closing
			if err == err {
				break
			}
			e.log.Fatal("fail to read p2p data", zap.Error(err))
		}

		if str == "" {
			return
		}
		if str != "\n" {

			mess := message{}
			if err := json.Unmarshal([]byte(str), &mess); err != nil {
				log.Fatal(err)
			}

			e.log.Info("New event read", zap.String("type", mess.Name.String()))

			switch mess.Name {
			case event.BlockChain:
				e.readBlockChain(mess.Value)
			case event.NewBlock:
				e.readNewBlock(mess.Value)
			case event.Wallet:
				e.readWallets(mess.Value)
			case event.Pool:
				e.readPool(mess.Value)
			case event.Files:
				e.readFilesAsk(mess.Value)
			}
		}
	}
}

// get blockChain for the first time
func (e *EndPoint) readBlockChain(chain []byte) {
	// si la blockChain est déjà renseignée, on ne fait rien
	if len(chain) <= len(blockchain.BlockChain) {
		e.log.Info("blockChain received smaller than current")
		return
	}

	mutex.Lock()

	var newBlockChain []blockchain.Block
	if err := json.Unmarshal(chain, &newBlockChain); err != nil {
		e.log.Error("fail to unmarshal blockChain received", zap.Error(err))
		return
	}

	if isValid := blockchain.IsValid(newBlockChain); !isValid {
		e.log.Fatal("fail to load blockChain")
		return
	}

	lastHashInDB, err := e.persistence.GetLastHash()
	if err != nil {
		e.log.Fatal("fail to get last hash in database", zap.Error(err))
		return
	}

	// on recherche le dernier hash trouvé dans la blockChain reçue
	j := e.getNumOnNewBlockChain(newBlockChain, lastHashInDB)
	for i := j; i < len(newBlockChain); i++ {
		current := newBlockChain[i]
		prevBlock := newBlockChain[i-1]
		serializeBLock, err := utils.Serialize(&current)
		if err != nil {
			e.log.Fatal("fail to Serialize block", zap.Error(err))
			return
		}
		// revérifie la cohérence des données reçues
		if blockchain.IsBlockValid(current, prevBlock) {
			// on met à jour la blockChain
			blockchain.BlockChain = append(blockchain.BlockChain, current)
			// on met à jour la BDD avec ces nouveaux hash reçus
			e.persistence.Update(current.Hash, serializeBLock)
		} else {
			e.log.Fatal("un block dans la blockchain reçue n'est pas valide")
			return
		}
	}

	mutex.Unlock()
	e.log.Info("received blockChain update")
}

func (e *EndPoint) getNumOnNewBlockChain(newBlockChain []blockchain.Block, lastHashInDB []byte) int {
	for i := len(newBlockChain); i > 0; i-- {
		if res := bytes.Compare(lastHashInDB, newBlockChain[i].Hash); res == 0 {
			return i
		}
	}
	return 0
}

func (e *EndPoint) readNewBlock(chain []byte) {
	var newBlock blockchain.Block

	if err := json.Unmarshal(chain, &newBlock); err != nil {
		e.log.Error("fail to unmarshal blockChain received", zap.Error(err))
		return
	}

	// test si le block est bien le suivant du block actuellement en base
	res := bytes.Compare(newBlock.PrevHash, blockchain.BlockChain[len(blockchain.BlockChain)-1].Hash)
	if res == 0 {

		if blockchain.IsBlockValid(newBlock, blockchain.BlockChain[len(blockchain.BlockChain)-1]) {
			mutex.Lock()
			blockchain.BlockChain = append(blockchain.BlockChain, newBlock)
			mutex.Unlock()

			ser, err := utils.Serialize(&newBlock)
			e.Handle(err)

			e.event.Push(event.BlockChain)

			err = e.persistence.Update(newBlock.Hash, ser)
			e.Handle(err)
			spew.Dump(blockchain.BlockChain)
		}

		// sinon, on test si c'est un ancien block, dans ce cas, on ne fait rien
	} else {
		for _, block := range blockchain.BlockChain {
			if res := bytes.Compare(block.Hash, newBlock.PrevHash); res == 0 {
				e.log.Error("cannot insert new block, it's an old hash")
				return
			}
		}
		e.log.Error("cannot insert new block, old hash not found")
		return
	}
}

func (e *EndPoint) readWallets(chain []byte) {
	if len(chain) <= len(e.wallets.Seeds) {
		e.log.Info("blockChain received smaller than current")
		return
	}

	mutex.Lock()
	if err := json.Unmarshal(chain, &e.wallets.Seeds); err != nil {
		e.log.Error("fail to unmarshal blockChain received", zap.Error(err))
		return
	}
	mutex.Unlock()
	e.log.Info("received blockChain update")
}

func (e *EndPoint) readPool(_ []byte) {

}

// on renotifie wallets and blockChain
func (e *EndPoint) readFilesAsk(_ []byte) {
	e.event.Push(event.BlockChain)
	e.event.Push(event.Wallet)
}
