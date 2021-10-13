package wallet

import (
	"bytes"
	"crypto/elliptic"
	"encoding/gob"
	"fmt"
	"github.com/ariden83/blockchain/config"
	"io/ioutil"
	"log"
	"os"
)

type Wallets struct {
	FilePath string
	Seeds    []Seed
}

type Options struct {
	// Required options.
	Dir      string
	ValueDir string
	File     string
	ReadOnly bool
	InMemory bool
}

func Init(conf *config.Config) (*Wallets, error) {
	wallets := Wallets{
		FilePath: conf.Wallet.File,
	}
	wallets.Seeds = make([]Seed, 0)

	err := wallets.LoadFile(conf.Wallet)

	return &wallets, err
}

func (Wallets) DefaultOptions(conf config.Wallet) Options {
	return Options{
		Dir:      conf.Path,
		ValueDir: conf.Path,
		File:     conf.File,
	}
}

// var _ io.Reader = (*os.File)(nil)

func (ws *Wallets) LoadFile(conf config.Wallet) error {
	opt := ws.DefaultOptions(conf)

	var wallets Wallets

	if _, err := os.Stat(conf.File); os.IsNotExist(err) {
		if err := ws.createDirs(opt); err != nil {
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

func (ws *Wallets) createDirs(opt Options) error {
	for _, path := range []string{opt.Dir, opt.ValueDir} {
		dirExists, err := ws.exists(path)
		if err != nil {
			return fmt.Errorf("Invalid Dir: %s error: %+v", path, err)
		}
		if !dirExists {
			if opt.ReadOnly {
				return fmt.Errorf("Cannot find directory %q for read-only open", path)
			}
			// Try to create the directory
			err = os.MkdirAll(path, 0700)
			if err != nil {
				return fmt.Errorf("Error Creating Dir: %s error: %+v", path, err)
			}
		}
	}
	return nil
}

func (Wallets) exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}
