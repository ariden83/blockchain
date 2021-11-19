package persistence

import (
	"github.com/ariden83/blockchain/config"
	"github.com/dgraph-io/badger"
	"os"
)

type Persistence struct {
	db       *badger.DB
	config   config.Database
	lastHash []byte
}

type IPersistence interface {
	GetLastHash() ([]byte, error)
	Update([]byte, []byte) error
	LastHash() []byte
	GetCurrentHashSerialize(hash []byte) ([]byte, error)
	DBExists() bool
	SetLastHash(lastHash []byte)
	Close()
}

// InitBlockChain will be what starts a new blockChain
func Init(conf config.Database) (*Persistence, error) {
	opts := badger.DefaultOptions(conf.Path)

	db, err := badger.Open(opts)
	if err != nil {
		return nil, err
	}

	per := &Persistence{
		db:     db,
		config: conf,
	}

	return per, nil
}

func (p *Persistence) DBExists() bool {
	if _, err := os.Stat(p.config.File); os.IsNotExist(err) {
		return false
	}
	return true
}

func (p *Persistence) LastHash() []byte {
	return p.lastHash
}

func (p *Persistence) SetLastHash(lastHash []byte) {
	p.lastHash = lastHash
}

func (p *Persistence) Update(lastHash []byte, hashSerialize []byte) error {
	err := p.db.Update(func(txn *badger.Txn) error {
		// "lh" stand for last hash
		err := txn.Set(lastHash, hashSerialize)
		if err != nil {
			return err
		}
		err = txn.Set([]byte("lh"), lastHash)
		p.lastHash = lastHash
		return err
	})
	return err
}

func (p *Persistence) GetLastHash() ([]byte, error) {
	var lastHash []byte

	err := p.db.View(func(txn *badger.Txn) error {
		// "lh" stand for last hash
		if _, err := txn.Get([]byte("lh")); err == badger.ErrKeyNotFound {
			return nil
		} else {
			item, err := txn.Get([]byte("lh"))
			if err != nil {
				return err
			}
			err = item.Value(func(val []byte) error {
				lastHash = val
				return nil
			})
			return err
		}
	})

	return lastHash, err
}

func (p *Persistence) GetCurrentHashSerialize(hash []byte) ([]byte, error) {
	var currentHashSerialize []byte
	err := p.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(hash)
		if err != nil {
			return err
		}
		err = item.Value(func(val []byte) error {
			currentHashSerialize = val
			return nil
		})
		return err
	})
	return currentHashSerialize, err
}

func (p *Persistence) Close() error {
	return p.db.Close()
}
