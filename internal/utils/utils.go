package utils

import (
	"bytes"
	"encoding/gob"
	"github.com/ariden83/blockchain/internal/blockchain"
)

func DeserializeBlock(data []byte) (*blockchain.Block, error) {
	var block blockchain.Block

	decoder := gob.NewDecoder(bytes.NewReader(data))

	err := decoder.Decode(&block)

	return &block, err
}

type SerializeInput interface{}

func Serialize(b SerializeInput) ([]byte, error) {
	var res bytes.Buffer
	encoder := gob.NewEncoder(&res)

	err := encoder.Encode(b)
	return res.Bytes(), err
}
