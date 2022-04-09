package p2p

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/ariden83/blockchain/internal/blockchain"
	"github.com/ariden83/blockchain/internal/event"
	"github.com/ariden83/blockchain/internal/iterator"
	"github.com/ariden83/blockchain/internal/p2p/address"
	"github.com/ariden83/blockchain/internal/p2p/validation"
	"github.com/davecgh/go-spew/spew"
	"github.com/satori/go.uuid"
	"go.uber.org/zap"
	"sync"
	"time"
)

var mutex = &sync.Mutex{}

// routine Go qui diffuse le dernier état de notre blockchain toutes les 5 secondes à nos pairs
// Ils le recevront et le jetteront si la longueur est plus courte que la leur. Ils l'accepteront si c'est plus long
func (e *EndPoint) writeData(rw *bufio.ReadWriter) {

	go func() {
		for {
			e.event.Push(event.Message{
				Type:  event.Address,
				Value: address.RecreateAddress(),
			})

			now := time.Now()
			nextTick := time.Date(now.Year(), now.Month(), now.Day(), now.Hour()+1, 0, 0, 0, time.Local)
			nextTick = nextTick.Add(e.cfg.AddressTimer)
			timer := nextTick.Sub(time.Now())

			time.Sleep(timer)

			e.log.Info("recreate address list")
		}

	}()

	go func() {
		var bytes []byte

		e.writerReady = true
		for mess := range e.event.NewReader() {
			e.log.Info("New event push", zap.String("type", mess.Type.String()), zap.String("ID", mess.ID))

			switch mess.Type {
			case event.BlockChain:
				bytes = e.sendBlockChain(rw)
			case event.BlockChainFull:
				bytes = e.sendBlockChainFull(rw)
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
			fmt.Println(fmt.Sprintf("********************************************* writeData %s", string(bytes)))
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
	fmt.Println(fmt.Sprintf("********************************************* marshal %s", string(allBytes)))
	rw.WriteString(fmt.Sprintf("%s\n", string(allBytes)))
	rw.Flush()
	mutex.Unlock()
}

type callFiles struct {
	Token string
}

func (e *EndPoint) callFiles(_ *bufio.ReadWriter) []byte {
	bytes, err := json.Marshal(callFiles{
		Token: e.cfg.Token,
	})
	if err != nil {
		e.log.Error("fail to marshal files message to send", zap.Error(err))
		return []byte{}
	}

	return bytes
}

func (e *EndPoint) sendAddress(_ *bufio.ReadWriter) []byte {

	bytes, err := json.Marshal(address.GetNewAddress())
	if err != nil {
		e.log.Error("fail to marshal new address", zap.Error(err))
		return []byte{}
	}

	return bytes
}

func (e *EndPoint) resendBlock(_ *bufio.ReadWriter) []byte {
	bytes, err := json.Marshal(blockchain.BlockChain)
	if err != nil {
		e.log.Error("fail to marshal blockChain", zap.Error(err))
		return []byte{}
	}

	return bytes
}

func (e *EndPoint) sendBlock(_ *bufio.ReadWriter, block validation.Validator) []byte {
	bytes, err := json.Marshal(block)
	if err != nil {
		e.log.Error("fail to marshal block message to send", zap.Error(err))
		return []byte{}
	}
	return bytes
}

func (e *EndPoint) sendBlockChain(_ *bufio.ReadWriter) []byte {
	spew.Dump(blockchain.BlockChain)
	bytes, err := json.Marshal(blockchain.BlockChain)
	if err != nil {
		e.log.Error("fail to marshal blockChain", zap.Error(err))
		return []byte{}
	}

	return bytes
}

func (e *EndPoint) sendBlockChainFull(_ *bufio.ReadWriter) []byte {

	blocks := []blockchain.Block{}
	iterator := iterator.New(e.persistence)
	for {
		block, err := iterator.Next()
		if err != nil {
			e.log.Error("fail to iterate next block", zap.Error(err))
		}
		blocks = append(blocks, *block)
		if len(block.PrevHash) == 0 {
			break
		}
	}

	bytes, err := json.Marshal(&blocks)
	if err != nil {
		e.log.Error("fail to marshal blockChain", zap.Error(err))
		return []byte{}
	}

	return bytes
}

// @todo envoyer en stream
func (e *EndPoint) sendWallets(_ *bufio.ReadWriter) []byte {
	listSeeds, err := e.wallets.GetSeeds()
	if err != nil {
		e.log.Error("fail to marshal wallets", zap.Error(err))
		return []byte{}
	}

	bytes, err := json.Marshal(listSeeds)
	if err != nil {
		e.log.Error("fail to marshal wallets", zap.Error(err))
		return []byte{}
	}

	return bytes
}

func (e *EndPoint) sendPool(_ *bufio.ReadWriter) []byte {
	return []byte{}
}
