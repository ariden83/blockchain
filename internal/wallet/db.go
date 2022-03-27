package wallet

import (
	"bytes"
	"os"

	"github.com/dgraph-io/badger"
	"go.uber.org/zap"

	"github.com/ariden83/blockchain/internal/utils"
	pkgError "github.com/ariden83/blockchain/pkg/errors"
)

func (w *Wallets) DBExists() bool {
	if _, err := os.Stat(w.File); os.IsNotExist(err) {
		return false
	}
	return true
}

func (w *Wallets) isSeedInDB(seed *Seed) error {
	if w.db != nil {
		var valCopy []byte

		if err := w.db.View(func(txn *badger.Txn) error {
			item, err := txn.Get(seed.Address)
			if err != nil {
				return err
			}

			if err := item.Value(func(val []byte) error {
				valCopy = append([]byte{}, val...)
				return nil
			}); err != nil {
				return err
			}
			return nil

		}); err != nil {
			w.log.Error("fail to get mnemonic in database", zap.Error(err))
			return pkgError.ErrorSeedPasswordInvalid
		}

		s := &SeedNoPrivKey{}
		if err := utils.Deserialize(valCopy, s); err != nil {
			w.log.Error("fail to deserialize mnemonic", zap.Error(err))
			return pkgError.ErrorSeedPasswordInvalid
		}
		seed.ExtraData = map[string]interface{}{
			"password": s.Password,
		}

	} else {
		for _, s := range w.Seeds {
			if res := bytes.Compare(s.Address, seed.Address); res == 0 {
				seed.ExtraData = map[string]interface{}{
					"password": s.Password,
				}
			}
		}
	}

	return nil
}

func (w *Wallets) saveInDB(seed SeedNoPrivKey) error {
	if w.db == nil {
		w.removeFromLocalMemory(seed)
		w.Seeds = append(w.Seeds, seed)
		return nil
	}
	serializeBLock, err := utils.Serialize(&seed)
	if err != nil {
		w.log.Error("cannot serialize new seed", zap.Error(err))
		return pkgError.ErrInternalDependencyError
	}
	return w.db.Update(func(txn *badger.Txn) error {
		e := badger.NewEntry(seed.Address, serializeBLock)
		return txn.SetEntry(e)
	})
}

func (w *Wallets) updateInDB(seed SeedNoPrivKey) error {
	if w.db == nil {
		w.Seeds = append(w.Seeds, seed)
		return nil
	}
	serializeBLock, err := utils.Serialize(&seed)
	if err != nil {
		w.log.Error("cannot serialize new seed", zap.Error(err))
		return pkgError.ErrInternalDependencyError
	}
	return w.db.Update(func(txn *badger.Txn) error {
		e := badger.NewEntry(seed.Address, serializeBLock)
		return txn.SetEntry(e)
	})
}

func (w *Wallets) removeFromLocalMemory(seed SeedNoPrivKey) {
	for i, s := range w.Seeds {
		if bytes.Compare(s.Address, seed.Address) == 0 {
			w.Seeds = append(w.Seeds[:i], w.Seeds[i+1:]...)
		}
	}
}
