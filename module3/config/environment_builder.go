package config

import (
	"log"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Server struct {
		Port        int           `env:"PORT"`
		Timeout     time.Duration `env:"TIMEOUT"`
		IdleTimeout time.Duration `env:"IDLE_TIMEOUT"`
	}
	Database struct {
		Server       string `env:"DB_SERVER"`
		Port         int    `env:"DB_PORT"`
		Username     string `env:"DB_USER"`
		Password     string `env:"DB_PASSWORD"`
		DatabaseName string `env:"DB_NAME"`
	}
	Authorization struct {
		Username string `env:"AUTH_USERNAME"`
		Password string `env:"AUTH_PASSWORD"`
	}
}

func LoadConfig() *Config {
	var cfg Config
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		log.Fatal("Error while loading environment variables: %v\n", err)
	}

	return &cfg
}
