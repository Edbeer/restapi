package config

import (
	"log"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

// Config
type Config struct {
	Server ServerConfig `yaml:"server"`
	Logger Logger       `yaml:"logger"`
}

// Server config struct
type ServerConfig struct {
	Port              string `yaml:"Port"`
	Mode              string `yaml:"Mode"`
	JwtSecretKey      string `yaml:"JwtSecretKey"`
	CookieName        string `yaml:"CookieName"`
	ReadTimeout       int    `yaml:"ReadTimeout"`
	WriteTimeout      int    `yaml:"WriteTimeout"`
	SSL               bool   `yaml:"SSL"`
	CtxDefaultTimeout int    `yaml:"CtxDefaultTimeout"`
}

type Logger struct {
	Development       bool
	DisableCaller     bool
	DisableStacktrace bool
	Encoding          string
	Level             string
}

var (
	config *Config
	once   sync.Once
)

// Get the config file
func GetConfig() *Config {
	once.Do(func() {
		log.Println("read application configuration")
		config = &Config{}
		if err := cleanenv.ReadConfig("config/config.yml", config); err != nil {
			help, _ := cleanenv.GetDescription(config, nil)
			log.Println(help)
			log.Fatal(err)
		}
	})
	return config
}
