package config

import (
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env         string `yaml:"env" env:"ENV" env-default:"local"`
	HTTPServer  `yaml:"http_server"`
	Database    `yaml:"database"`
	JWT         `yaml:"jwt"`
}

type HTTPServer struct {
	Address string `yaml:"address" env:"SERVER_ADDRESS" env-default:":8080"`
}

type Database struct {
	Host     string `yaml:"host" env:"DB_HOST" env-default:"localhost"`
	Port     int    `yaml:"port" env:"DB_PORT" env-default:"5432"`
	User     string `yaml:"user" env:"DB_USER" env-default:"postgres"`
	Password string `yaml:"password" env:"DB_PASSWORD"`
	Name     string `yaml:"name" env:"DB_NAME" env-default:"warehouse"`
	SSLMode  string `yaml:"sslmode" env:"DB_SSLMODE" env-default:"disable"`
}

type JWT struct {
	Secret     string `yaml:"secret" env:"JWT_SECRET"`
	Expiration int    `yaml:"expiration" env:"JWT_EXPIRATION_HOURS" env-default:"24"`
}

// MustLoad loads configuration from config file and environment variables
// Panics if loading fails
func MustLoad(configPath string) *Config {
	// Check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	var cfg Config

	// Read config from YAML file
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("failed to read config: %s", err)
	}

	// Override with environment variables
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		log.Fatalf("failed to read env: %s", err)
	}

	return &cfg
}
