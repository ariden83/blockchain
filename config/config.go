package config

import (
	"context"
	"fmt"
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
	Path     string `config:"wallet_path"`
	File     string `config:"wallet_file"`
	WithFile bool   `config:"wallet_with_file"`
}

type Metrics struct {
	Namespace string
	Name      string
	Port      int
}

type Transactions struct {
	Reward *big.Int
}

// address of actual user which mine on this server
type Miner struct {
	PubKey     string `config:"miner_pub_key"`
	PrivateKey string `config:"miner_private_key"`
	Address    string `config:"miner_address"`
}

type API struct {
	Enabled bool `config:"api_enabled"`
	Port    int  `config:"api_port"`
}

type P2P struct {
	Enabled bool `config:"p2p_enabled"`
	// Parse options from the command line
	// Port ouvre le port auquel nous voulons autoriser les connexions
	Port int `config:"p2p_port"`
	// secio : sécurisation des flux
	Secio bool `config:"p2p_secio_enabled"`
	// target nous permet de spécifier l'adresse d'un autre hôte auquel nous voulons nous connecter,
	// ce qui signifie que nous agissons en tant qu'homologue d'un hôte si nous utilisons ce drapeau.
	Target string `config:"p2p_target"`
	// seed est le paramètre aléatoire facultatif utilisé pour construire notre adresse
	// que d'autres pairs peuvent utiliser pour se connecter à nous
	Seed int64 `config:"p2p_seed"`

	TimeToCommunicate int `config:"p2p_time_to_communicate"`
	// token utilisé pour assurer la sécurité de la connexion
	Token string
}

type Log struct {
	Path     string `config:"log_path"`
	WithFile bool   `config:log_with_file"`
}

type Config struct {
	Name    string
	Version string
	Address string
	Threads int `json:"threads"`
	//reward is the amnount of tokens given to someone that "mines" a new block
	Gas          Gas
	Wallet       Wallet
	Metrics      Metrics
	Transactions Transactions
	Miner        Miner
	Database     Database
	API          API
	Log          Log
	P2P          P2P
}

func getDefaultConfig() *Config {
	return &Config{
		Name: "blockChain",
		Transactions: Transactions{
			Reward: big.NewInt(100),
		},
		Database: Database{
			Path: "./tmp/blocks",
			File: "./tmp/blocks/MANIFEST",
		},
		Wallet: Wallet{
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
		Miner: Miner{},
		API: API{
			Enabled: true,
			Port:    8098,
		},
		P2P: P2P{
			Port:              8097,
			Enabled:           true,
			TimeToCommunicate: 5,
			Target:            "",
		},
		Log: Log{
			Path: "./tmp/logs",
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
