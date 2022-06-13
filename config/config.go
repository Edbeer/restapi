package config

import (
	"log"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

// Config
type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Logger   Logger         `yaml:"logger"`
	Postgres PostgresConfig `yaml:"postgres"`
	Redis    RedisConfig    `yaml:"redis"`
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

// Postgresql config
type PostgresConfig struct {
	PostgresqlHost     string `yaml:"PostgresqlHost"`
	PostgresqlPort     string `yaml:"PostgresqlPort"`
	PostgresqlUser     string `yaml:"PostgresqlUser"`
	PostgresqlPassword string `yaml:"PostgresqlPassword"`
	PostgresqlDbname   string `yaml:"PostgresqlDbname"`
	PostgresqlSSLMode  bool   `yaml:"PostgresqlSSLMode"`
	PgDriver           string `yaml:"PgDriver"`
}

type RedisConfig struct {
	RedisAddr      string `yaml:"RedisAddr"`
	RedisPassword  string `yaml:"RedisPassword"`
	RedisDB        string `yaml:"RedisDB"`
	RedisDefaultdb string `yaml:"RedisDefaultdb"`
	MinIdleConns   int    `yaml:"MinIdleConns"`
	PoolSize       int    `yaml:"PoolSize"`
	PoolTimeout    int    `yaml:"PoolTimeout"`
	Password       string `yaml:"Password"`
	DB             int    `yaml:"DB"`
}

type Logger struct {
	Development       bool   `yaml:"Development"`
	DisableCaller     bool   `yaml:"DisableCaller"`
	DisableStacktrace bool   `yaml:"DisableStacktrace"`
	Encoding          string `yaml:"Encoding"`
	Level             string `yaml:"Level"`
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
