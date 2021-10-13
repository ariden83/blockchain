package config

import (
	"context"
	"fmt"
	"github.com/ariden83/blockchain/internal/handle"
	"github.com/heetch/confita"
	"github.com/heetch/confita/backend"
	"github.com/heetch/confita/backend/env"
	"github.com/heetch/confita/backend/flags"
	"math/big"
	"reflect"
	"strings"
)

type Database struct {
	Path string `config:"blockchain_db_path"`
	File string `config:"blockchain_file"`
}

type Gas struct {
	// Minimum amount of Wei per GAS to be paid for a transaction to be accepted for mining. Note: gasprice is listed in wei. Note 2: --gasprice is a “Legacy Option”
	Price int `config:"gas-price"`
	// Les mineurs sur le réseau décident de la limite de gaz du bloc.
	// Il faut réintroduire une fonction de limite de gaz adaptative
	Limit int `config:"gas-limit"`
	// il existe une stratégie d'exploitation minière par défaut
	// d'une limite de gaz de bloc minimale de 4 712 388 pour la plupart des clients
	TargetMinLimit int `config:"targetgasMinlimit"`
	TargetMaxLimit int `config:"targetgaslimit"`
	// Amount of gas per block to target when sealing a new block
	FloorTarget int `config:"gas-floor-target"`
	// A cap on how large we will raise the gas limit per block due to transaction volume
	Cap int `config:"gas-cap"`
}

type Wallet struct {
	Path string `config:"wallet_path"`
	File string `config:"wallet_file"`
}

type Metrics struct {
	Namespace string
	Name      string
	Port      int
}

type Config struct {
	Name     string
	Version  string
	Port     int
	Database Database
	Address  string
	//reward is the amnount of tokens given to someone that "mines" a new block
	Reward  *big.Int
	Gas     Gas
	Wallet  Wallet
	Metrics Metrics
}

func getDefaultConfig() *Config {
	return &Config{
		Name:   "blockChain",
		Port:   8098,
		Reward: big.NewInt(100),
		Database: Database{
			Path: "./tmp/blocks",
			File: "./tmp/blocks/MANIFEST",
		},
		Wallet: Wallet{
			Path: "./tmp/wallets",
			File: "./tmp/wallets/wallets.data",
		},
		Gas: Gas{
			Price:          4000000000,
			Limit:          1,
			TargetMinLimit: 1,
			TargetMaxLimit: 4712388,
			FloorTarget:    4700000,
			Cap:            6283184,
		},
		Metrics: Metrics{
			Namespace: "block",
			Name:      "chain",
		},
	}
}

// New Load the config
func New() *Config {
	loaders := []backend.Backend{
		env.NewBackend(),
		flags.NewBackend(),
	}

	loader := confita.NewLoader(loaders...)

	cfg := getDefaultConfig()
	err := loader.Load(context.Background(), cfg)
	if err != nil {
		handle.Handle(err)
	}

	fmt.Println(fmt.Sprintf("%+v", cfg))
	return cfg
}

func (c *Config) String() string {
	val := reflect.ValueOf(c).Elem()
	s := "\n-------------------------------\n"
	s += "-  Application configuration  -\n"
	s += "-------------------------------\n"
	for i := 0; i < val.NumField(); i++ {
		v := val.Field(i)
		t := val.Type().Field(i)
		c.applyWithType(&s, "", v, t)
	}
	return s
}

func (c *Config) applyWithType(s *string, parentKey string, v reflect.Value, k reflect.StructField) {
	obfuscate := false

	tag := k.Tag.Get("config")
	if idx := strings.Index(tag, ","); idx != -1 {
		opts := strings.Split(tag[idx+1:], ",")

		for _, opt := range opts {
			if opt == "obfuscate" {
				obfuscate = true
			}
		}
	}
	if !obfuscate {
		if parentKey != "" {
			parentKey += "-"
		}
		switch v.Kind() {
		case reflect.String:
			*s += fmt.Sprintf("%s: \"%v\"\n", parentKey+k.Name, v.Interface())
			return
		case reflect.Bool:
		case reflect.Int:
			*s += fmt.Sprintf("%s: %v\n", parentKey+k.Name, v.Interface())
			return
		case reflect.Struct:
			parentKey += k.Name
			c.DeepStructFields(s, parentKey, v.Interface())
			return
		}
	}
}

func (c *Config) DeepStructFields(s *string, parentKey string, iface interface{}) {
	ifv := reflect.ValueOf(iface)
	ift := reflect.TypeOf(iface)

	for i := 0; i < ift.NumField(); i++ {
		v := ifv.Field(i)
		t := ift.Field(i)
		c.applyWithType(s, parentKey, v, t)
	}
}
