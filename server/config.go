package server

import "github.com/jinzhu/configor"

// Config - loaded from config.yml
// Defaults are for local development with sqlite3
type Config struct {
	Database struct {
		Driver     string `default:"sqlite3"`
		Connection string `default:"../tmp/dev.db"`
	}
	App struct {
		Addr   string `default:"127.0.0.1"`
		Port   string `default:"3000"`
		Secret string `default:"my-secret-key"`
	}
}

// LoadConfig ...
func LoadConfig() *Config {
	config := &Config{}
	configor.Load(&config, "../config.yml")
	return config
}
