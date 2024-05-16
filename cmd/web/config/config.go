package config

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/heetch/confita"
	"github.com/heetch/confita/backend"
	"github.com/heetch/confita/backend/env"
	"github.com/heetch/confita/backend/file"
	"github.com/heetch/confita/backend/flags"

	"github.com/ariden83/blockchain/config"
)

type OAuthConfig struct {
	ClientStore  string
	ClientID     string
	ClientSecret string
	Scopes       []string
	URLAPI       string
	Enable       bool
}

type Auth struct {
	GoogleAPI OAuthConfig
	Classic   OAuthConfig
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

type Metadata struct {
	Title string `config:"metadata_title"`
}

type ReCaptcha struct {
	SiteKey   string        `config:"recaptcha_sitekey"`
	SecretKey string        `config:"recaptcha_secretkey"`
	Timeout   time.Duration `config:"recaptcha_timeout"`
	URL       string
}

type Locales struct {
	Path string `config:"locales_path"`
	Lang []string
}

type Mails struct {
	PublicKey string `config:"mails_publickey"`
	SecretKey string `config:"mails_secretkey"`
	ProxyURL  string `config:"mails_proxyurl"`
}

type Config struct {
	Name          string `config:"name" yaml:"name"`
	Version       string `config:"version"`
	DumpVar       bool   `config:"dump_var"`
	Domain        string `config:"domain"`
	TemplatesDir  string `config:"template_dir" yaml:"template_dir"`
	StaticDir     string `config:"static_dir" yaml:"static_dir"`
	Port          int    `config:"port"`
	Log           config.Log
	Api           Api
	Auth          Auth
	Metrics       Metrics
	Healthz       Healthz
	BlockchainAPI BlockchainAPI
	Metadata      Metadata
	ReCaptcha     ReCaptcha
	Locales       Locales
	Mails         Mails
}

func getDefaultConfig() *Config {
	return &Config{
		Name:         "blockChain",
		Version:      "0.0.0",
		Port:         4000,
		DumpVar:      false,
		Domain:       "http://localhost:4000",
		TemplatesDir: "cmd/web/templates/",
		StaticDir:    "./cmd/web/static/",
		Log: config.Log{
			CLILevel: "info",
			WithFile: false,
		},
		Api: Api{
			Port: 8098,
		},
		Auth: Auth{
			Classic: OAuthConfig{
				ClientID:     uuid.New().String()[:8],
				ClientSecret: uuid.New().String()[:8],
				Scopes:       []string{"all"},
				Enable:       true,
			},
			GoogleAPI: OAuthConfig{
				ClientID:     "",
				ClientSecret: "",
				URLAPI:       "https://www.googleapis.com/oauth2/v2/userinfo?access_token=",
				Enable:       false,
			},
		},
		Mails: Mails{},
		Metrics: Metrics{
			Port: 8101,
			Host: "0.0.0.0",
		},
		BlockchainAPI: BlockchainAPI{
			MaxSizeCall: 1024 * 1024 * 12,
			URL:         "0.0.0.0:8155",
			TimeOut:     10,
		},
		Metadata: Metadata{
			Title: "blockchain-altcoin",
		},
		ReCaptcha: ReCaptcha{
			SiteKey:   "",
			SecretKey: "",
			URL:       "https://www.google.com/recaptcha/api/siteverify",
		},
		Locales: Locales{
			Path: "./cmd/web/locales/",
			Lang: []string{"en-US", "fr-FR"},
		},
	}
}

// New Load the config
func New() (*Config, error) {
	loaders := []backend.Backend{
		env.NewBackend(),
		flags.NewBackend(),
	}

	configFile := findConfigFilePathRecursively("prod", 0)
	if configFile != "" {
		loaders = append(loaders, file.NewBackend(configFile))
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

func findConfigFilePathRecursively(environment string, depth int) string {
	char := "../"
	if depth == 0 {
		char = "./"
	}
	if depth > 3 {
		return ""
	}

	filePath := strings.Repeat(char, depth) + "cmd/web/config/config." + environment + ".yaml"
	if _, err := os.Stat(filePath); err == nil {
		return filePath
	}
	depth++

	return findConfigFilePathRecursively(environment, depth)
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
