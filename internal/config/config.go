package config

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"

	"github.com/ARUMANDESU/go-revise/pkg/env"
)

type Config struct {
	EnvMode         env.Mode      `yaml:"env_mode"         env:"ENV_Mode"         env-default:"local"`
	StartTimeout    time.Duration `yaml:"start_timeout"    env:"START_TIMEOUT"    env-default:"10s"`
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout" env:"SHUTDOWN_TIMEOUT" env-default:"10s"`
	Telegram        Telegram      `yaml:"telegram"`
	HTTP            HTTP          `yaml:"http"`
	DatabaseURL     string        `yaml:"database_url"     env:"DATABASE_URL"`
}

type Telegram struct {
	EnvMode    env.Mode `yaml:"-"           env:"-"`
	Token      string   `yaml:"token"       env:"TELEGRAM_TOKEN"`
	WebhookURL string   `yaml:"webhook_url" env:"TELEGRAM_WEBHOOK_URL"`
	URL        string   `yaml:"url"         env:"TELEGRAM_URL"`
}

type HTTP struct {
	Port string `yaml:"port" env:"HTTP_PORT" env-default:"8080"`
}

func MustLoad() Config {
	path := fetchConfigPath()
	if path == "" {
		return mustLoadFromEnv()
	}

	return mustLoadByPath(path)
}

func mustLoadByPath(configPath string) Config {
	cfg, err := loadByPath(configPath)
	if err != nil {
		panic(err)
	}

	return cfg
}

func loadByPath(configPath string) (Config, error) {
	var cfg Config

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return Config{}, fmt.Errorf("there is no config file: %w", err)
	}

	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		return Config{}, fmt.Errorf("failed to read config: %w", err)
	}

	cfg.Telegram.EnvMode = cfg.EnvMode
	return cfg, nil
}

func mustLoadFromEnv() Config {
	var cfg Config

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		panic("Mode empty")
	}

	cfg.Telegram.EnvMode = cfg.EnvMode
	return cfg
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
