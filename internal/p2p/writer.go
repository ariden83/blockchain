package p2p

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/ariden83/blockchain/internal/blockchain"
	"github.com/ariden83/blockchain/internal/event"
	"github.com/satori/go.uuid"
	"go.uber.org/zap"
	"sync"
)

var mutex = &sync.Mutex{}

// routine Go qui diffuse le dernier état de notre blockchain toutes les 5 secondes à nos pairs
// Ils le recevront et le jetteront si la longueur est plus courte que la leur. Ils l'accepteront si c'est plus long
func (e *EndPoint) writeData(rw *bufio.ReadWriter) {
	go func() {
		var bytes []byte

		e.writerReady = true
		for mess := range e.event.NewReader() {
			e.log.Info("New event push", zap.String("type", mess.Type.String()), zap.String("ID", mess.ID))

			switch mess.Type {
			case event.BlockChain:
				bytes = e.sendBlockChain(rw)
			case event.Wallet:
				bytes = e.sendWallets(rw)
			case event.Pool:
				bytes = e.sendPool(rw)
			case event.Files:
				bytes = e.callFiles(rw)
			case event.Address:
				bytes = e.sendAddress(rw)
			case event.NewBlock:
				// resend message only
				bytes = mess.Value
			default:
				e.log.Error(fmt.Sprintf("Event type push not found %+v", mess))
			}

			e.marshal(rw, mess, bytes)
		}

	}()

	go func() {
		for block := range e.event.NewBlockReader() {
			bytes := e.sendBlock(rw, block)
			e.marshal(rw, event.Message{Type: event.NewBlock}, bytes)
		}
	}()
}

func (e *EndPoint) marshal(rw *bufio.ReadWriter, evt event.Message, bytes []byte) {
	if len(bytes) == 0 {
		return
	}

	mess := event.Message{
		Type:  evt.Type,
		Value: bytes,
		ID:    evt.ID,
	}

	if mess.ID == "" {
		mess.ID = uuid.NewV4().String()
	}

	e.saveMsgReceived(mess.ID)

	allBytes, err := json.Marshal(mess)
	if err != nil {
		e.log.Error("fail to marshal message", zap.Error(err))
		return
	}

	mutex.Lock()
	rw.WriteString(fmt.Sprintf("%s\n", string(allBytes)))
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
		e.log.Error("fail to marshal files message to send", zap.Error(err))
		return []byte{}
	}

	return bytes
}

func (e *EndPoint) sendAddress(rw *bufio.ReadWriter) []byte {
	return []byte{}
}

func (e *EndPoint) resendBlock(rw *bufio.ReadWriter) []byte {
	bytes, err := json.Marshal(blockchain.BlockChain)
	if err != nil {
		e.log.Error("fail to marshal blockChain", zap.Error(err))
		return []byte{}
	}

	return bytes
}

func (e *EndPoint) sendBlock(rw *bufio.ReadWriter, block blockchain.Block) []byte {
	bytes, err := json.Marshal(block)
	if err != nil {
		e.log.Error("fail to marshal block message to send", zap.Error(err))
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
