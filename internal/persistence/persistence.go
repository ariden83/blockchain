package persistence

import (
	"github.com/ariden83/blockchain/internal/handle"
	"github.com/dgraph-io/badger"
)

type Persistence struct {
	db       *badger.DB
	LastHash []byte
}

// InitBlockChain will be what starts a new blockChain
func Init(dbPath string) *Persistence {
	opts := badger.DefaultOptions(dbPath)

	db, err := badger.Open(opts)
	handle.Handle(err)

	p := &Persistence{
		db: db,
	}

	return p
}

func (p *Persistence) SetLastHash(lastHash []byte) {
	p.LastHash = lastHash
}

func (p *Persistence) Update(lastHash []byte, hashSerialize []byte) error {
	err := p.db.Update(func(txn *badger.Txn) error {
		// "lh" stand for last hash
		err := txn.Set(lastHash, hashSerialize)
		if err != nil {
			return err
		}
		err = txn.Set([]byte("lh"), lastHash)
		p.LastHash = lastHash
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

func (p *Persistence) Close() {
	p.db.Close()
}