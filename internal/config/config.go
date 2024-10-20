package config

import (
	"errors"
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"time"
)

const (
	configPathEnv = "LINKER_CONFIG_PATH"
)

type Config struct {
	AvailableMethods []string         `yaml:"available_methods"`
	Client           HttpClientConfig `yaml:"client"`
}

type HttpClientConfig struct {
	BasePath string        `yaml:"base_path"`
	Timeout  time.Duration `yaml:"timeout"`
}

type GrpcClientConfig struct {
}

func Load() (*Config, error) {
	configPath := os.Getenv(configPathEnv)
	if configPath == "" {
		return nil, errors.New("переменная LINKER_CONFIG_PATH не установлена")
	}

	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		return nil, fmt.Errorf("неудалось загрузить конфиг: %w", err)
	}

	return &cfg, nil
}
