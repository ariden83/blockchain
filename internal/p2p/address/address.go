package address

import (
	"encoding/json"
	"sync"
)

var IAM Address

func init() {
	IAM = New()
}

type Address struct {
	// CurrentAddress is only use by API call.
	currentAddress []string
	// NewAddress is used to share address on network.
	newAddress []string
	// IAM represent my address.
	iam   string
	mutex *sync.Mutex
}

func New() Address {
	return Address{
		mutex: &sync.Mutex{},
	}
}

func (a *Address) Address() string {
	return a.iam
}

func (a *Address) SetAddress(myAddress string) {
	a.iam = myAddress
	a.newAddress = []string{a.iam}
	a.currentAddress = a.newAddress
}

func (a *Address) RecreateMyAddress() []byte {
	a.currentAddress = a.newAddress
	a.newAddress = []string{a.iam}

	bytes, err := json.Marshal(a.newAddress)
	if err != nil {
		return []byte{}
	}
	return bytes
}

func (a *Address) CurrentAddress() []string {
	return a.currentAddress
}

func (a *Address) GetNewAddress() []string {
	return a.newAddress
}

func (a *Address) UpdateAddress(received []string) []string {
	toNewAddress := map[string]bool{}
	for _, addr := range received {
		toNewAddress[addr] = true
	}
	for _, addr := range a.newAddress {
		toNewAddress[addr] = true
	}

	saveAddress := []string{}
	for key := range toNewAddress {
		saveAddress = append(saveAddress, key)
	}
	a.mutex.Lock()
	a.newAddress = saveAddress
	a.mutex.Unlock()
	return a.newAddress
}
