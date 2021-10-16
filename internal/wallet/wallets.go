package wallet

import (
	"bytes"
	"crypto/elliptic"
	"encoding/gob"
	"github.com/ariden83/blockchain/config"
	"github.com/ariden83/blockchain/internal/dir"
	"io/ioutil"
	"log"
	"os"
)

type Wallets struct {
	FilePath string
	Seeds    []Seed
}

func Init(conf config.Wallet) (*Wallets, error) {
	wallets := Wallets{
		FilePath: conf.File,
	}
	wallets.Seeds = make([]Seed, 0)

	var err error
	if conf.WithFile {
		err = wallets.LoadFile(conf)
	}
	return &wallets, err
}

// var _ io.Reader = (*os.File)(nil)

func (ws *Wallets) LoadFile(conf config.Wallet) error {
	opt := dir.Options{
		Dir:      conf.Path,
		File:     conf.File,
		FileMode: os.FileMode(0700),
	}

	var wallets Wallets

	if !dir.DirExists(opt) {
		if err := dir.CreateDir(opt); err != nil {
			return err
		}
		ws.Save()
	}

	fileContent, err := ioutil.ReadFile(conf.File)
	if err != nil {
		return err
	}

	gob.Register(elliptic.P256())
	decoder := gob.NewDecoder(bytes.NewReader(fileContent))
	err = decoder.Decode(&wallets)
	if err != nil {
		return err
	}

	ws.Seeds = wallets.Seeds

	return nil
}

func (ws *Wallets) Save() {
	var content bytes.Buffer

	gob.Register(elliptic.P256)

	encoder := gob.NewEncoder(&content)
	err := encoder.Encode(ws)
	if err != nil {
		log.Panic(err)
	}

	err = ioutil.WriteFile(ws.FilePath, content.Bytes(), 0644)
	if err != nil {
		log.Panic(err)
	}
}
