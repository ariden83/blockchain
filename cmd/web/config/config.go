package config

import (
	"context"
	"fmt"
	"github.com/ariden83/blockchain/config"
	"github.com/heetch/confita"
	"github.com/heetch/confita/backend"
	"github.com/heetch/confita/backend/env"
	"github.com/heetch/confita/backend/flags"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"
)

type Token struct {
	SecretKey string `config:"token_secret_key,obfuscate"`
}

type Api struct {
	Port    int `config:"api_port"`
	TimeOut float64
}

type BlockchainAPI struct {
	MaxSizeCall int     `config:"blockchainapi_maxcallsize"`
	URL         string  `config:"blockchainapi_url"`
	TimeOut     float64 `config:"blockchainapi_timeout"`
}

type Metrics struct {
	Port int    `config:"metrics_port"`
	Host string `config:"metrics_host"`
}

type Healthz struct {
	ReadTimeout  time.Duration `config:"healthz_read_timeout"`
	WriteTimeout time.Duration `config:"healthz_write_timeout"`
}

type Config struct {
	Name          string `config:"name"`
	Version       string `config:"version"`
	TemplatesDir  string `config:"template_dir"`
	StaticDir     string `config:"static_dir"`
	StaticRoute   string `config:"static_route"`
	Port          int    `config:"port"`
	Log           config.Log
	Api           Api
	Token         Token
	Metrics       Metrics
	Healthz       Healthz
	BlockchainAPI BlockchainAPI
}

func getDefaultConfig() *Config {
	return &Config{
		Name:         "blockChain",
		Version:      "0.0.0",
		Port:         4000,
		TemplatesDir: "cmd/web/templates/",
		StaticDir:    "./cmd/web/static/",
		StaticRoute:  "/static/",
		Log: config.Log{
			CLILevel: "info",
			WithFile: false,
		},
		Api: Api{
			Port: 8098,
		},
		Token: Token{
			SecretKey: "chihuahua",
		},
		Metrics: Metrics{
			Port: 8101,
			Host: "0.0.0.0",
		},
		BlockchainAPI: BlockchainAPI{
			MaxSizeCall: 1024 * 1024 * 12,
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

	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)
	fmt.Println(fmt.Sprintf("Current path %s", exPath))

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

// BuildPort buils a port string from a port number.
func (c *Config) BuildPort(port int) string {
	return fmt.Sprintf(":%d", port)
}

func (c *Config) GetExplorerPort() int {
	return c.Port
}
