package config

import (
	"context"
	"fmt"
	"math/big"
	"reflect"
	"strings"
	"time"

	"github.com/heetch/confita"
	"github.com/heetch/confita/backend"
	"github.com/heetch/confita/backend/env"
	"github.com/heetch/confita/backend/flags"

	"github.com/ariden83/blockchain/internal/event/trace"
	p2pfactory "github.com/ariden83/blockchain/internal/p2p/factory"
	p2padapter "github.com/ariden83/blockchain/internal/p2p/impl/p2p"
	persistencefactory "github.com/ariden83/blockchain/internal/persistence/factory"
	"github.com/ariden83/blockchain/internal/persistence/impl/badger"
	"github.com/ariden83/blockchain/internal/wallet"
)

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

type Metrics struct {
	Namespace string
	Name      string
	Port      int
}

type Transactions struct {
	Reward *big.Int
	Miner
}

type Miner struct {
	PubKey     string `config:"miner_pub_key"`
	PrivateKey string `config:"miner_private_key"`
	Address    string `config:"miner_address"`
}

type API struct {
	Enabled bool   `config:"api_enabled"`
	Port    int    `config:"api_port"`
	Host    string `config:"api_host" yaml:"api_host"`
}

type GRPC struct {
	Enabled bool   `config:"grpc_enabled"`
	Port    int    `config:"grpc_port" yaml:"grpc_port"`
	Host    string `config:"grpc_host" yaml:"grpc_host"`
}

type Log struct {
	Path     string `config:"log_path"`
	WithFile bool   `config:"log_with_file"`
	CLILevel string `config:"log_cli_level" yaml:"log_cli_level"`
}

type Config struct {
	Name    string
	Version string
	Address string
	Threads int `json:"threads"`
	//reward is the amnount of tokens given to someone that "mines" a new block
	Gas          Gas
	Wallet       wallet.Config
	Metrics      Metrics
	Transactions Transactions
	Database     persistencefactory.Config
	API          API
	Log          Log
	P2P          p2pfactory.Config
	XCache       p2padapter.XCache
	GRPC         GRPC
	Trace        trace.Config
}

func getDefaultConfig() *Config {
	return &Config{
		Name:    "blockChain",
		Version: "0.0.0",
		Transactions: Transactions{
			Reward: big.NewInt(100),
			Miner: Miner{
				PubKey: "xpub661MyMwAqRbcFTZYiEcSv4Qj2Qr2NzQ7rjYc3iv9c6VSTxoYsqA9AA6nNbp8e9nVR9hRARXz5CApP6j5BxUnohyj89oSg3zZdDuKmGhdSFF",
			},
		},
		Trace: trace.Config{
			Enabled: true,
		},
		Database: persistencefactory.Config{
			Implementation: persistencefactory.ImplementationBadger,
			Badger: badger.Config{
				Path: "./tmp/blocks",
				File: "./tmp/blocks/MANIFEST",
			},
		},
		Wallet: wallet.Config{
			Path:     "./tmp/wallets",
			File:     "./tmp/wallets/wallets.data",
			WithFile: true,
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
			Port:      8099,
		},
		API: API{
			Enabled: true,
			Port:    8098,
			Host:    "0.0.0.0",
		},
		GRPC: GRPC{
			Enabled: true,
			Port:    8155,
			Host:    "0.0.0.0",
		},
		P2P: p2pfactory.Config{
			Implementation: p2pfactory.ImplementationP2P,
			Config: p2padapter.Config{
				Port:               8097,
				Enabled:            true,
				TimeToCommunicate:  5,
				ProtocolID:         "/p2p/1.0.0",
				DiscoveryNamespace: "blockchain",
				AddressTimer:       30 * time.Minute,
			},
		},
		XCache: p2padapter.XCache{
			Size:            5000,
			TTL:             60,
			MaxSizeAccepted: 60000,
			NegSize:         500,
			NegTTL:          30,
			Active:          true,
		},
		Log: Log{
			Path:     "./tmp/logs",
			CLILevel: "info",
		},
	}
}

// New Load the config
func New() (*Config, error) {
	loaders := []backend.Backend{
		env.NewBackend(),
		flags.NewBackend(),
	}

	loader := confita.NewLoader(loaders...)

	cfg := getDefaultConfig()
	err := loader.Load(context.Background(), cfg)
	if err != nil {
		return cfg, err
	}

	fmt.Println(fmt.Sprintf("%+v", cfg))
	return cfg, nil
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
