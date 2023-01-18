// Package badger represent a persistence storage system service, based on badger lib.
package badger

import (
	"os"

	"github.com/dgraph-io/badger"
)

// Persistence represent a persistence adapter structure.
type Persistence struct {
	db       *badger.DB
	config   Config
	lastHash []byte
}

// Config represent a persistence badger config.
type Config struct {
	Path string
	File string
}

// New represent a new persistence adapter.
func New(conf Config) (*Persistence, error) {
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

// DBExists test if a DB already exist.
func (p *Persistence) DBExists() bool {
	if _, err := os.Stat(p.config.File); os.IsNotExist(err) {
		return false
	}
	return true
}

// LastHash get the last hash linked to the current DB.
func (p *Persistence) LastHash() []byte {
	return p.lastHash
}

// SetLastHash update the last hash linked to the current DB.
func (p *Persistence) SetLastHash(lastHash []byte) {
	p.lastHash = lastHash
}

// Update the last hash into the current DB.
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

// GetLastHash the last hash in the current DB.
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
				lastHash = append([]byte{}, val...)
				return nil
			})
			return err
		}
	})

	return lastHash, err
}

// GetCurrentHashSerialize the last hash serialize in the current DB.
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

// Close the current DB.
func (p *Persistence) Close() error {
	return p.db.Close()
}
