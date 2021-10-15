package p2p

import (
	"bufio"
	"encoding/json"
	"github.com/ariden83/blockchain/internal/blockchain"
	"github.com/ariden83/blockchain/internal/event"
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

func (e *EndPoint) readBlockChain(chain []byte) {
	if len(chain) <= len(blockchain.BlockChain) {
		e.log.Info("blockChain received smaller than current")
		return
	}
	mutex.Lock()

	if err := json.Unmarshal(chain, &blockchain.BlockChain); err != nil {
		e.log.Error("fail to unmarshal blockChain received", zap.Error(err))
		return
	}
	mutex.Unlock()
	e.log.Info("received blockChain update")
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
