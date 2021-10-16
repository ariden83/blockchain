package p2p

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/ariden83/blockchain/internal/blockchain"
	"github.com/ariden83/blockchain/internal/event"
	"go.uber.org/zap"
)

// routine Go qui diffuse le dernier état de notre blockchain toutes les 5 secondes à nos pairs
// Ils le recevront et le jetteront si la longueur est plus courte que la leur. Ils l'accepteront si c'est plus long
func (e *EndPoint) writeData(rw *bufio.ReadWriter) {
	go func() {
		var bytes []byte

		for data := range e.event.NewReader() {
			e.log.Info("New event push", zap.String("type", data.String()))
			mutex.Lock()

			switch data {
			case event.BlockChain:
				bytes = e.sendBlockChain(rw)
			case event.Wallet:
				bytes = e.sendWallets(rw)
			case event.Pool:
				bytes = e.sendPool(rw)
			case event.Files:
				bytes = e.callFiles(rw)
			}
			mutex.Unlock()

			e.marshal(rw, data, bytes)
		}
	}()

	go func() {
		for block := range e.event.NewBlockReader() {
			e.log.Info("New block push")
			mutex.Lock()
			bytes := e.sendBlock(rw, block)
			mutex.Unlock()

			e.marshal(rw, event.NewBlock, bytes)
		}
	}()
}

func (e *EndPoint) marshal(rw *bufio.ReadWriter, evt event.EventType, bytes []byte) {
	if len(bytes) == 0 {
		return
	}

	mess := message{
		Name:  evt,
		Value: bytes,
	}

	bytes, err := json.Marshal(mess)
	if err != nil {
		e.log.Error("fail to marshal message", zap.Error(err))
		return
	}

	mutex.Lock()
	rw.WriteString(fmt.Sprintf("%s\n", string(bytes)))
	rw.Flush()
	mutex.Unlock()
}

type callFiles struct {
	Token string
}

func (e *EndPoint) callFiles(rw *bufio.ReadWriter) []byte {
	bytes, err := json.Marshal(callFiles{
		Token: e.cfg.Token,
	})
	if err != nil {
		e.log.Error("fail to marshal blockChain", zap.Error(err))
		return []byte{}
	}

	return bytes
}

func (e *EndPoint) sendBlock(rw *bufio.ReadWriter, block blockchain.Block) []byte {
	bytes, err := json.Marshal(block)
	if err != nil {
		e.log.Error("fail to marshal blockChain", zap.Error(err))
		return []byte{}
	}

	return bytes
}

func (e *EndPoint) sendBlockChain(rw *bufio.ReadWriter) []byte {
	bytes, err := json.Marshal(blockchain.BlockChain)
	if err != nil {
		e.log.Error("fail to marshal blockChain", zap.Error(err))
		return []byte{}
	}

	return bytes
}

func (e *EndPoint) sendWallets(rw *bufio.ReadWriter) []byte {
	bytes, err := json.Marshal(e.wallets.GetSeeds())
	if err != nil {
		e.log.Error("fail to marshal wallets", zap.Error(err))
		return []byte{}
	}

	return bytes
}

func (e *EndPoint) sendPool(rw *bufio.ReadWriter) []byte {
	return []byte{}
}
