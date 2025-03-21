package config

import (
	"errors"
	"log/slog"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server ServerConfig `yaml:"server"`
}

type ServerConfig struct {
	Host        string        `yaml:"host"`
	Port        string        `yaml:"port"`
	Timeout     time.Duration `yaml:"timeout"`
	IdleTimeout time.Duration `yaml:"idle_timeout"`
}

var (
	Cfg = &Config{}
)

func LoadConfig() {
	configPath := "./config/config.yaml"
	f, err := os.Open(configPath)
	if errors.Is(err, os.ErrNotExist) {
		slog.Error("config file does not exist", "error", err)
		os.Exit(1)
	} else if err != nil {
		slog.Error("config file error", "error", err)
		os.Exit(1)
	}

	if err = yaml.NewDecoder(f).Decode(Cfg); err != nil {
		slog.Error("cannot read config", "error", err)
		os.Exit(1)
	}
}
