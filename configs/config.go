package configs

import (
	"fmt"
	"path"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		Server   ServerConfig
		Database DatabaseConfig
		JWT      JWTConfig
	}

	ServerConfig struct {
		Port string `env-required:"true" yaml:"port" env:"HTTP_PORT"`
	}

	DatabaseConfig struct {
		Host     string `env-required:"false" yaml:"host" env:"DB_HOST"`
		Port     string `env-required:"false" yaml:"port" env:"DB_PORT"`
		Name     string `env-required:"false" yaml:"name" env:"DB_NAME"`
		User     string `env-required:"true" env:"DB_USER"`
		Password string `env-required:"true" env:"DB_PASSWORD"`
		SSLMode  string `env-required:"true" env:"SSL_MODE"`
	}

	JWTConfig struct {
		Secret string        `env-required:"true" env:"JWT_SECRET"`
		TTL    time.Duration `env-required:"false" yaml:"name" env:"TTL"`
	}
)

func LoadConfig(configPath string) (*Config, error) {
	cfg := &Config{}

	err := cleanenv.ReadConfig(path.Join("./", configPath), cfg)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	err = cleanenv.UpdateEnv(cfg)
	if err != nil {
		return nil, fmt.Errorf("error updating env: %w", err)
	}

	return cfg, nil
}
