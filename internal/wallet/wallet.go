package wallet

import (
	"bytes"
	"encoding/gob"
	"errors"
	"go.uber.org/zap"
	"os"
	"sync"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/LuisAcerv/btchdwallet/crypt"
	"github.com/brianium/mnemonic"
	"github.com/dgraph-io/badger"
	"github.com/wemeetagain/go-hdwallet"
	"golang.org/x/crypto/sha3"

	"github.com/ariden83/blockchain/config"
	"github.com/ariden83/blockchain/internal/utils"
	pkgError "github.com/ariden83/blockchain/pkg/errors"
)

type Wallets struct {
	File      string
	Seeds     []Seed
	TempSeeds []Seed
	withFile  bool
	db        *badger.DB
	log       *zap.Logger
}

type IWallets interface {
	Create([]byte) (*Seed, error)
	Close() error
	DBExists() bool
	GetAllPublicSeeds() []SeedNoPrivKey
	GetSeed([]byte, []byte) (*SeedNoPrivKey, error)
	GetSeeds() *[]Seed
	UpdateSeeds([]Seed)
	Validate(pubKey []byte) bool
}

var ErrorSeedPasswordInvalid = errors.New("invalid seed password")

// Seed represents each 'item' in the blockchain
type Seed struct {
	Address   string
	Timestamp int64
	PubKey    string
	PrivKey   string
	Mnemonic  []byte
	ExtraData map[string]interface{}
}

func (s *Seed) verifyPassword(password []byte) error {
	if err := bcrypt.CompareHashAndPassword(s.ExtraData["password"].([]byte), password); err != nil {
		return ErrorSeedPasswordInvalid
	}
	return nil
}

type SeedNoPrivKey struct {
	Timestamp int64
	PubKey    string
	Address   string
}

var (
	mutex = &sync.Mutex{}
)

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

func (w *Wallets) DBExists() bool {
	if _, err := os.Stat(w.File); os.IsNotExist(err) {
		return false
	}
	return true
}

func (w *Wallets) GetAllPublicSeeds() []SeedNoPrivKey {
	var allSeeds []SeedNoPrivKey
	for _, j := range w.Seeds {
		allSeeds = append(allSeeds, SeedNoPrivKey{
			Timestamp: j.Timestamp,
			Address:   j.Address,
		})
	}
	return allSeeds
}

// @todo update wallet database
func (w *Wallets) UpdateSeeds(seed []Seed) {
	w.Seeds = seed
}

func (w *Wallets) Save(seed Seed) error {
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
		e := badger.NewEntry([]byte(seed.Mnemonic), serializeBLock)
		return txn.SetEntry(e)
	})
}

func (w *Wallets) Create(password []byte) (*Seed, error) {
	password, err := encryptPassword(password)
	if err != nil {
		w.log.Error("fail to encrypt password", zap.Error(err))
		return nil, pkgError.ErrInternalDependencyError
	}

	seed := crypt.CreateHash()
	mnemonic, err := mnemonic.New([]byte(seed), mnemonic.English)
	if err != nil {
		w.log.Error("fail to generate new mnemonic", zap.Error(err))
		return nil, pkgError.ErrInternalDependencyError
	}

	// Create a master private key
	masterPrv := hdwallet.MasterKey([]byte(mnemonic.Sentence()))

	// Convert a private key to public key
	masterPub := masterPrv.Pub()

	// Get your address
	address := masterPub.Address()

	t := time.Now().UnixNano() / int64(time.Millisecond)

	mnemonicStr := mnemonic.Sentence()
	mnemonicHash := hash([]byte(mnemonicStr))
	newSeed := Seed{
		Address:   address,
		PubKey:    masterPub.String(),
		PrivKey:   masterPrv.String(),
		Mnemonic:  mnemonicHash,
		Timestamp: t,
		ExtraData: map[string]interface{}{
			"password": password,
		},
	}

	w.TempSeeds = append(w.TempSeeds, newSeed)

	return &Seed{
		Mnemonic:  []byte(mnemonicStr),
		Address:   newSeed.Address,
		PubKey:    newSeed.PubKey,
		Timestamp: newSeed.Timestamp,
	}, nil
}

func (w *Wallets) Validate(pubKey []byte) bool {
	for i, s := range w.TempSeeds {
		if s.PubKey == string(pubKey) {
			w.Save(s)
			w.TempSeeds = append(w.TempSeeds[:i], w.TempSeeds[i+1:]...)
			return true
		}
	}
	return false
}

func (w *Wallets) GetSeed(mnemonic, password []byte) (*SeedNoPrivKey, error) {
	var (
		err  error
		seed *Seed
	)

	mnemonicHash := hash(mnemonic)
	if w.db != nil {
		var valCopy []byte
		if err := w.db.View(func(txn *badger.Txn) error {
			item, err := txn.Get(mnemonicHash)
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
			return nil, pkgError.ErrSeedNotFound
		}
		seed, err = deserialize(valCopy)
		if err != nil {
			w.log.Error("fail to deserialize mnemonic", zap.Error(err))
			return nil, pkgError.ErrInternalDependencyError
		}

	} else {
		for _, s := range w.Seeds {
			if res := bytes.Compare(s.Mnemonic, mnemonicHash); res == 0 {
				seed = &s
			}
		}
	}

	if seed == nil {
		return nil, pkgError.ErrSeedNotFound
	}

	if err := seed.verifyPassword(password); err != nil {
		w.log.Info("invalid password", zap.Error(err))
		return nil, pkgError.ErrInvalidPassword
	}

	return &SeedNoPrivKey{
		PubKey:  seed.PubKey,
		Address: seed.Address,
	}, nil
}

func (w *Wallets) GetSeeds() *[]Seed {
	return &w.Seeds
}

func (w *Wallets) Close() error {
	if w.db != nil {
		return w.db.Close()
	}
	return nil
}

func GetPubKey(privKey []byte) ([]byte, error) {
	masterPrv, err := hdwallet.StringWallet(string(privKey))
	if err != nil {
		return []byte{}, err
	}
	return []byte(masterPrv.Pub().String()), nil
}

func deserialize(data []byte) (*Seed, error) {
	var seed Seed
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&seed)
	return &seed, err
}

func encryptPassword(password []byte) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword(password, bcrypt.MinCost)
	if err != nil {
		return []byte{}, errors.New("fail to encrypt password")
	}
	return hash, nil
}

const localKey = "blockchain"

func hash(mnemonic []byte) []byte {
	fixedSlice := sha3.Sum512(append(mnemonic, []byte(localKey)...))
	byteSlice := make([]byte, 64)
	for i := 0; i < len(fixedSlice); i++ {
		byteSlice[i] = fixedSlice[i]
	}
	return byteSlice
}
