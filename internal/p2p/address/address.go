package address

import (
	"encoding/json"
	"sync"
)

var (
	// only use by API call
	CurrentAddress []string
	// use to share address on network
	NewAddress []string
	IAM        string
	mutex      = &sync.Mutex{}
)

func GetMe() string {
	return IAM
}

func SetIAM(iam string) {
	IAM = iam
	NewAddress = []string{iam}
	CurrentAddress = NewAddress
}

func RecreateAddress() []byte {
	CurrentAddress = NewAddress
	NewAddress = []string{IAM}

	bytes, err := json.Marshal(NewAddress)
	if err != nil {
		return []byte{}
	}
	return bytes
}

func GetCurrentAddress() []string {
	return CurrentAddress
}

func GetNewAddress() []string {
	return NewAddress
}

func UpdateAddress(received []string) []string {
	toNewAddress := map[string]bool{}
	for _, addr := range received {
		toNewAddress[addr] = true
	}
	for _, addr := range NewAddress {
		toNewAddress[addr] = true
	}

	saveAddress := []string{}
	for key := range toNewAddress {
		saveAddress = append(saveAddress, key)
	}
	mutex.Lock()
	NewAddress = saveAddress
	mutex.Unlock()
	return NewAddress
}
