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

func Serialize(b SerializeInput) ([]byte, error) {
	var res bytes.Buffer
	encoder := gob.NewEncoder(&res)

	err := encoder.Encode(b)
	return res.Bytes(), err
}

func RandomString(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	s := make([]rune, n)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}
