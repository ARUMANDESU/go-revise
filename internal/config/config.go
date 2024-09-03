package config

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env             string        `yaml:"env" env:"ENV" env-default:"local"`
	StartTimeout    time.Duration `yaml:"start_timeout" env:"START_TIMEOUT" env-default:"10s"`
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout" env:"SHUTDOWN_TIMEOUT" env-default:"10s"`
	Telegram        Telegram      `yaml:"telegram"`
	DatabaseURL     string        `yaml:"database_url" env:"DATABASE_URL"`
}

type Telegram struct {
	Token      string `yaml:"token" env:"TELEGRAM_TOKEN"`
	WebhookURL string `yaml:"webhook_url" env:"TELEGRAM_WEBHOOK_URL"`
	URL        string `yaml:"url" env:"TELEGRAM_URL"`
}

func MustLoad() *Config {
	path := fetchConfigPath()
	if path == "" {
		return mustLoadFromEnv()
	}

	return mustLoadByPath(path)
}

func mustLoadByPath(configPath string) *Config {
	cfg, err := loadByPath(configPath)
	if err != nil {
		panic(err)
	}

	return cfg
}

func loadByPath(configPath string) (*Config, error) {
	var cfg Config

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("there is no config file: %w", err)
	}

	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	return &cfg, nil
}

func mustLoadFromEnv() *Config {
	var cfg Config

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		panic("Env empty")
	}
	return &cfg
}

func fetchConfigPath() string {
	var res string

	if flag.Lookup("config") == nil {
		flag.StringVar(&res, "config", "", "path to config file")
		flag.Parse()
	}

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}
