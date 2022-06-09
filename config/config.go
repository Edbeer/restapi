package config

// Config
type Config struct {
	Server ServerConfig `yaml:"server"`
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
