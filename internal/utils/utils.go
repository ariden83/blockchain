package utils

import (
	"bytes"
	"encoding/gob"
	"math/rand"
)

type DeserializedOutput interface{}

func Deserialize(data []byte, i DeserializedOutput) error {
	decoder := gob.NewDecoder(bytes.NewReader(data))

	err := decoder.Decode(i)

	return err
}

type SerializeInput interface{}

func Serialize(serializeStr SerializeInput) ([]byte, error) {
	var res bytes.Buffer
	encoder := gob.NewEncoder(&res)

	err := encoder.Encode(serializeStr)
	return res.Bytes(), err
}

func RandomString(n uint8) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	s := make([]rune, n)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}
