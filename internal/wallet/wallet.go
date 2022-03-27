package wallet

import (
	"bytes"
	"sync"

	"go.uber.org/zap"

	"github.com/LuisAcerv/btchdwallet/crypt"
	"github.com/brianium/mnemonic"
	"github.com/dgraph-io/badger"

	"github.com/ariden83/blockchain/config"
	"github.com/ariden83/blockchain/internal/utils"
	pkgError "github.com/ariden83/blockchain/pkg/errors"
)

var (
	mutex = &sync.Mutex{}
)

type Wallets struct {
	File      string
	Seeds     []SeedNoPrivKey
	TempSeeds []SeedNoPrivKey
	withFile  bool
	db        *badger.DB
	log       *zap.Logger
}

type IWallets interface {
	Create([]byte) (*Seed, error)
	Close() error
	DBExists() bool
	GetSeed([]byte, []byte) (*Seed, error)
	GetSeeds() []SeedNoPrivKey
	UpdateSeeds([]SeedNoPrivKey)
	Validate([]byte) bool
}

func Init(cfg config.Wallet, log *zap.Logger) (*Wallets, error) {
	var err error
	opts := badger.DefaultOptions(cfg.Path)

	wallets := Wallets{
		File:     cfg.File,
		withFile: cfg.WithFile,
		log:      log.With(zap.String("service", "wallet")),
	}
	if cfg.WithFile {
		if wallets.db, err = badger.Open(opts); err != nil {
			return nil, err
		}
	}
	return &wallets, err
}

func (w *Wallets) GetSeeds() []SeedNoPrivKey {
	var allSeeds []SeedNoPrivKey

	if w.db == nil {
		for _, j := range w.Seeds {
			allSeeds = append(allSeeds, SeedNoPrivKey{
				Timestamp: j.Timestamp,
				Address:   j.Address,
			})
		}
		return allSeeds
	}

	if err := w.db.View(func(txn *badger.Txn) error {

		// create a Badger iterator with the default settings
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10
		it := txn.NewIterator(opts)
		defer it.Close()

		// have the iterator walk the LMB tree
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			//k := item.Key()
			seed := &SeedNoPrivKey{}
			if err := item.Value(func(v []byte) (err error) {
				return utils.Deserialize(v, seed)
			}); err != nil {
				return err
			}

			allSeeds = append(allSeeds, *seed)
		}

		return nil
	}); err != nil {
		panic(err)
	}

	return allSeeds
}

// @todo update wallet database
func (w *Wallets) UpdateSeeds(seed []SeedNoPrivKey) {
	w.Seeds = seed
}

func (w *Wallets) Create(password []byte) (*Seed, error) {
	password, err := encryptPassword(password)
	if err != nil {
		w.log.Error("fail to encrypt password", zap.Error(err))
		return nil, pkgError.ErrInternalDependencyError
	}

	hash := crypt.CreateHash()
	mnemonic, err := mnemonic.New([]byte(hash), mnemonic.English)
	if err != nil {
		w.log.Error("fail to generate new mnemonic", zap.Error(err))
		return nil, pkgError.ErrInternalDependencyError
	}

	seed := w.allKeysFromMnemonic([]byte(mnemonic.Sentence()))
	seed.ExtraData = map[string]interface{}{
		"password": password,
	}

	w.TempSeeds = append(w.TempSeeds, SeedNoPrivKey{
		Address:   seed.Address,
		Password:  password,
		Timestamp: seed.Timestamp,
	})

	return seed, nil
}

func (w *Wallets) Validate(privKey []byte) bool {
	seed, err := w.allKeysFromPrivate(privKey)
	if err != nil {
		w.log.Error("fail to validate seed with private key", zap.Error(err))
		return false
	}

	for i, s := range w.TempSeeds {
		if bytes.Compare(s.Address, seed.Address) == 0 {
			w.saveInDB(s)
			w.TempSeeds = append(w.TempSeeds[:i], w.TempSeeds[i+1:]...)
			return true
		}
	}
	return false
}

func (w *Wallets) GetSeed(mnemonic, password []byte) (*Seed, error) {
	seed := w.allKeysFromMnemonic(mnemonic)
	encryptPassword, err := encryptPassword(password)
	if err != nil {
		w.log.Error("fail to encrypt password", zap.Error(err))
		return nil, pkgError.ErrInternalDependencyError
	}
	// verify if seed in database
	if err := w.isSeedInDB(seed); err != nil {
		w.log.Info("seed not in database", zap.Error(err))
		w.saveInDB(SeedNoPrivKey{
			Address:   seed.Address,
			Password:  encryptPassword,
			Timestamp: seed.Timestamp,
		})
		// verify seed password
	} else if err := seed.verifyPassword(password); err != nil {
		w.log.Info("invalid seed password", zap.Error(err))
		w.saveInDB(SeedNoPrivKey{
			Address:   seed.Address,
			Password:  encryptPassword,
			Timestamp: seed.Timestamp,
		})
	}

	return seed, nil
}

func (w *Wallets) Close() error {
	if w.db != nil {
		return w.db.Close()
	}
	return nil
}
