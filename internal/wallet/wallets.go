package wallet

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/LuisAcerv/btchdwallet/crypt"
	"github.com/brianium/mnemonic"
	"github.com/dgraph-io/badger"
	"github.com/wemeetagain/go-hdwallet"

	"github.com/ariden83/blockchain/config"
	"github.com/ariden83/blockchain/internal/utils"
)

type Wallets struct {
	File     string
	Seeds    []Seed
	withFile bool
	db       *badger.DB
}

type IWallets interface {
	Create([]byte) (*SeedNoPrivKey, error)
	Close() error
	DBExists() bool
	GetAllPublicSeeds() []SeedNoPrivKey
	GetSeed([]byte, []byte) (*SeedNoPrivKey, error)
	GetSeeds() *[]Seed
	UpdateSeeds([]Seed)
}

var ErrorSeedPasswordInvalid = errors.New("invalid seed password")

// Seed represents each 'item' in the blockchain
type Seed struct {
	Address   string
	Timestamp int64
	PubKey    string
	PrivKey   string
	Mnemonic  string
	Password  []byte
}

func (s *Seed) verifyPassword(password []byte) error {
	if err := bcrypt.CompareHashAndPassword(s.Password, password); err != nil {
		return ErrorSeedPasswordInvalid
	}
	return nil
}

type SeedNoPrivKey struct {
	Timestamp int64
	PubKey    string
	Mnemonic  string
	Address   string
}

var mutex = &sync.Mutex{}

var ErrorFailEncryptPassword = errors.New("fail to encrypt password")

var ErrorCannotSerialize = errors.New("cannot serialize new seed")

func Init(cfg config.Wallet) (*Wallets, error) {
	var err error
	opts := badger.DefaultOptions(cfg.Path)

	wallets := Wallets{
		File:     cfg.File,
		withFile: cfg.WithFile,
	}

	if wallets.db, err = badger.Open(opts); err != nil {
		return nil, err
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
			Mnemonic:  j.Mnemonic,
			Timestamp: j.Timestamp,
			PubKey:    j.PubKey,
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
	serializeBLock, err := utils.Serialize(&seed)
	if err != nil {
		return fmt.Errorf("%s %w", ErrorCannotSerialize, err)
	}
	return w.db.Update(func(txn *badger.Txn) error {
		e := badger.NewEntry([]byte(seed.Mnemonic), serializeBLock)
		return txn.SetEntry(e)
	})
}

func (w *Wallets) Create(password []byte) (*SeedNoPrivKey, error) {
	password, err := encryptPassword(password)
	if err != nil {
		return nil, err
	}

	seed := crypt.CreateHash()
	mnemonic, err := mnemonic.New([]byte(seed), mnemonic.English)
	if err != nil {
		return nil, err
	}

	// Create a master private key
	masterPrv := hdwallet.MasterKey([]byte(mnemonic.Sentence()))

	// Convert a private key to public key
	masterPub := masterPrv.Pub()

	// Get your address
	address := masterPub.Address()

	t := time.Now().UnixNano() / int64(time.Millisecond)
	newSeed := Seed{
		Address:   address,
		PubKey:    masterPub.String(),
		PrivKey:   masterPrv.String(),
		Mnemonic:  mnemonic.Sentence(),
		Password:  password,
		Timestamp: t,
	}

	w.Save(newSeed)

	return &SeedNoPrivKey{
		Mnemonic: newSeed.Mnemonic,
		Address:  newSeed.Address,
		PubKey:   newSeed.PubKey,
	}, nil
}

func (w *Wallets) GetSeed(mnemonic, password []byte) (*SeedNoPrivKey, error) {
	var valCopy []byte
	if err := w.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(mnemonic)
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
		return nil, err
	}

	seed, err := deserialize(valCopy)
	if err != nil {
		return nil, err
	}
	if err := seed.verifyPassword(password); err != nil {
		return nil, err
	}
	return &SeedNoPrivKey{
		PubKey:   seed.PubKey,
		Address:  seed.Address,
		Mnemonic: seed.Mnemonic,
	}, nil
}

func (w *Wallets) GetSeeds() *[]Seed {
	return &w.Seeds
}

func (w *Wallets) Close() error {
	return w.db.Close()
}

func deserialize(data []byte) (*Seed, error) {
	var seed Seed
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&seed)
	return &seed, err
}

func encryptPassword(password []byte) ([]byte, error) {
	// Generate "hash" to store from user password
	hash, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		return []byte{}, ErrorFailEncryptPassword
	}
	return hash, nil
}
