package utils

import (
	"bytes"
	"encoding/gob"
)

type DeserializedOutput interface{}

func Deserialize(data []byte, i DeserializedOutput) error {
	decoder := gob.NewDecoder(bytes.NewReader(data))

	err := decoder.Decode(i)

	return err
}

type SerializeInput interface{}

func Serialize(b SerializeInput) ([]byte, error) {
	var res bytes.Buffer
	encoder := gob.NewEncoder(&res)

	err := encoder.Encode(b)
	return res.Bytes(), err
}
