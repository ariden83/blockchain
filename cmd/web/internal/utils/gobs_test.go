package utils

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"reflect"
	"testing"
)

func TestToGob(t *testing.T) {
	type testStruct struct{ Test string }
	test := testStruct{Test: "I am a test struct."}

	t.Run("Should return a slice of bytes", func(t *testing.T) {
		result := ToGob(test)
		kind := reflect.TypeOf(result).Kind()

		if kind != reflect.Slice {
			t.Errorf("Should return a slice of bytes, got %s", kind)
		}
	})

	t.Run("Should return a decodable gob", func(t *testing.T) {
		encoded := ToGob(test)
		decoded := testStruct{}
		err := gob.NewDecoder(bytes.NewReader(encoded)).Decode(&decoded)

		if err != nil {
			t.Error("Unable to decode Gob.")
		}

		if decoded != test {
			t.Error("Decoded gob doesn't match original object.")
		}
	})

}

func ExampleToGob() {
	test := struct{ Test string }{Test: "I am a test struct."}
	encoded := ToGob(test)
	decoded := struct{ Test string }{}
	gob.NewDecoder(bytes.NewReader(encoded)).Decode(&decoded)
	fmt.Println(decoded)
	// Output: {I am a test struct.}
}

func TestFromGob(t *testing.T) {
	type testStruct struct{ Test string }
	test := testStruct{Test: "I am a test struct."}
	encoded := []byte{
		33, 255, 131, 3, 1, 1, 10, 116, 101, 115, 116, 83, 116,
		114, 117, 99, 116, 1, 255, 132, 0, 1, 1, 1, 4, 84, 101, 115, 116, 1, 12,
		0, 0, 0, 24, 255, 132, 1, 19, 73, 32, 97, 109, 32, 97, 32, 116, 101,
		115, 116, 32, 115, 116, 114, 117, 99, 116, 46, 0}

	t.Run("Should decode the gob into the original object", func(t *testing.T) {
		decoded := testStruct{}
		FromGob(&decoded, encoded)

		if !reflect.DeepEqual(test, decoded) {
			t.Errorf("Should return a %s, got %s", reflect.TypeOf(test), reflect.TypeOf(decoded))
		}

		if decoded != test {
			t.Error("Decoded gob doesn't match original object.")
		}
	})

}

func ExampleFromGob() {
	test := struct{ Test string }{Test: "I am a test struct."}
	var buffer bytes.Buffer
	gob.NewEncoder(&buffer).Encode(test)
	encoded := buffer.Bytes()
	decoded := struct{ Test string }{}
	FromGob(&decoded, encoded)
	fmt.Println(decoded)
	// Output: {I am a test struct.}
}
