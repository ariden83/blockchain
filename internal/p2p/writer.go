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

			if len(bytes) == 0 {
				continue
			}

			mess := message{
				Name:  data,
				Value: bytes,
			}

			bytes, err := json.Marshal(mess)
			if err != nil {
				e.log.Error("fail to marshal message", zap.Error(err))
				continue
			}

			mutex.Lock()
			rw.WriteString(fmt.Sprintf("%s\n", string(bytes)))
			rw.Flush()
			mutex.Unlock()
		}
	}()
}

type callFiles struct {
	Token string
}

func (e *EndPoint) callFiles(rw *bufio.ReadWriter) []byte {
	bytes, err := json.Marshal(callFiles{
		Token: e.cfg.P2P.Token,
	})
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
	bytes, err := json.Marshal(e.wallets.Seeds)
	if err != nil {
		e.log.Error("fail to marshal wallets", zap.Error(err))
		return []byte{}
	}

	return bytes
}

func (e *EndPoint) sendPool(rw *bufio.ReadWriter) []byte {
	return []byte{}
}
