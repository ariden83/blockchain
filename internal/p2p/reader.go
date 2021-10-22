package p2p

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ariden83/blockchain/internal/blockchain"
	"github.com/ariden83/blockchain/internal/event"
	"github.com/ariden83/blockchain/internal/utils"
	// "github.com/davecgh/go-spew/spew"
	"go.uber.org/zap"
)

var (
	uniqueMsg = "once"
)

func (e *EndPoint) saveMsgReceived(uid string) {
	e.msgReceived = append(e.msgReceived, uid)
}

func (e *EndPoint) msgAlreadyReceived(uid string) bool {
	if uid == uniqueMsg {
		return false
	}

	for _, a := range e.msgReceived {
		if a == uid {
			return true
		}
	}
	return false
}

// routine Go qui récupère le dernier état de notre blockchain toutes les 5 secondes
// err = rw.Flush()
func (e *EndPoint) readData(rw *bufio.ReadWriter) {
	go func() {
		e.readerReady = true
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

				mess := event.Message{}
				if err := json.Unmarshal([]byte(str), &mess); err != nil {
					e.log.Fatal("fail to unmarshal message received", zap.Error(err))
				}

				if e.msgAlreadyReceived(mess.ID) {
					continue
				}
				// save message ID received
				e.saveMsgReceived(mess.ID)

				e.log.Info("New event read", zap.String("type", mess.Type.String()))
				//spew.Dump(mess)
				switch mess.Type {
				case event.BlockChain:
					e.readBlockChain(mess.Value)
				case event.NewBlock:
					e.readNewBlock(mess)
				case event.Wallet:
					e.readWallets(mess.Value)
				case event.Pool:
					e.readPool(mess.Value)
				case event.Files:
					e.readFilesAsk(mess)
				default:
					e.log.Error(fmt.Sprintf("Event type received not found %+v", mess))
				}

				if mess.Type != event.Files {
					e.event.Push(mess)
				}
			}
		}
	}()
}

// get blockChain for the first time
func (e *EndPoint) readBlockChain(chain []byte) {

	var newBlockChain []blockchain.Block
	if err := json.Unmarshal(chain, &newBlockChain); err != nil {
		e.log.Error("fail to unmarshal blockChain received", zap.Error(err))
		return
	}

	// si la blockChain est déjà renseignée, on ne fait rien
	if len(newBlockChain) <= len(blockchain.BlockChain) {
		e.log.Info("blockChain received smaller or equal than current")
		return
	}

	// surcharge tout si blockchain actuelle est vide
	if len(blockchain.BlockChain) == 0 {
		blockchain.BlockChain = newBlockChain
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

	// on recherche le nombre de blocks chez nous non trouvé dans la blockChain reçue
	j := e.getNumOnNewBlockChain(newBlockChain, lastHashInDB)

	// à partir de ce numéro, on itere sur les nouveau blocks reçus pour les ajouter
	for i := len(newBlockChain) - j; i < len(newBlockChain); i++ {
		// si genesis (normalement on passe pas ici, sauf dans la version light)
		if i-1 < 0 {
			continue
		}
		current := newBlockChain[i]

		prevBlock := blockchain.GetLastBlock()
		serializeBLock, err := utils.Serialize(&current)
		if err != nil {
			e.log.Fatal("fail to Serialize block", zap.Error(err))
			return
		}

		// revérifie la cohérence des données reçues
		if blockchain.IsBlockValid(current, prevBlock) {
			e.log.Info("received blockChain update")
			// on met à jour la blockChain
			blockchain.BlockChain = append(blockchain.BlockChain, current)
			// on met à jour la BDD avec ces nouveaux hash reçus
			e.persistence.Update(current.Hash, serializeBLock)
		} else {
			e.log.Error("un block dans la blockchain reçue n'est pas valide")
			return
		}
	}
}

func (e *EndPoint) getNumOnNewBlockChain(newBlockChain []blockchain.Block, lastHashInDB []byte) int {
	if len(newBlockChain) == 0 {
		return 0
	}
	for i := range newBlockChain {
		if res := bytes.Compare(lastHashInDB, newBlockChain[len(newBlockChain)-1-i].Hash); res == 0 {
			return i
		}
	}
	return 0
}

func (e *EndPoint) readNewBlock(msg event.Message) {
	var (
		newBlock blockchain.Block
	)
	if err := json.Unmarshal(msg.Value, &newBlock); err != nil {
		e.log.Error("fail to unmarshal block received", zap.Error(err))
		//spew.Dump(string(msg.Value))
		return
	}

	if len(blockchain.BlockChain) == 0 {
		// on demande a recevoir la blockchain
		e.event.Push(event.Message{Type: event.BlockChain})
		return
	}

	/*fmt.Println(fmt.Sprintf("****************************************** BlockChain %d", len(blockchain.BlockChain)-1))
	spew.Dump(blockchain.BlockChain)
	fmt.Println("****************************************** block")
	spew.Dump(blockchain.BlockChain[len(blockchain.BlockChain)-1])
	fmt.Println("****************************************** newBlock")
	spew.Dump(newBlock)*/

	// test si le block est bien le suivant du block actuellement en base
	res := bytes.Compare(newBlock.PrevHash, blockchain.BlockChain[len(blockchain.BlockChain)-1].Hash)
	if res == 0 {

		if blockchain.IsBlockValid(newBlock, blockchain.BlockChain[len(blockchain.BlockChain)-1]) {

			ser, err := utils.Serialize(&newBlock)
			e.Handle(err)

			blockchain.BlockChain = append(blockchain.BlockChain, newBlock)
			err = e.persistence.Update(newBlock.Hash, ser)
			e.Handle(err)

			e.event.Push(event.Message{Type: event.BlockChain})
			//spew.Dump(blockchain.BlockChain)
		}

		// sinon, on test si c'est un ancien block, dans ce cas, on ne fait rien
	} else {
		for _, block := range blockchain.BlockChain {
			if res := bytes.Compare(block.Hash, newBlock.PrevHash); res == 0 {
				e.event.Push(event.Message{
					Type: event.BlockChain,
					ID:   msg.ID + "link",
				})
				e.log.Error("cannot insert new block, it's an old hash")
				return
			}
		}
		e.log.Error("cannot insert new block, old hash not found")
		return
	}
}

func (e *EndPoint) readWallets(chain []byte) {
	if len(chain) <= len(*e.wallets.GetSeeds()) {
		e.log.Info("blockChain received smaller than current")
		return
	}

	if err := json.Unmarshal(chain, e.wallets.GetSeeds()); err != nil {
		e.log.Error("fail to unmarshal blockChain received", zap.Error(err))
		return
	}
	e.log.Info("received wallets update")
	//spew.Dump(e.wallets.GetAllPublicSeeds())
}

func (e *EndPoint) readPool(_ []byte) {
	e.log.Info("************************************************************ readPool")
}

// on renotifie wallets and blockChain
func (e *EndPoint) readFilesAsk(m event.Message) {
	e.event.Push(event.Message{Type: event.BlockChain})
	e.event.Push(event.Message{Type: event.Wallet})
}
