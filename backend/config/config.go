package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type (
	// Config -.
	Config struct {
		App  `yaml:"app"`
		HTTP `yaml:"http"`
		PG   `yaml:"postgres"`
		S3   `yaml:"s3"`
	}

	// App -.
	App struct {
		Name    string `env-required:"true" yaml:"name"    env:"APP_NAME"`
		Version string `env-required:"true" yaml:"version" env:"APP_VERSION"`
	}

	// HTTP -.
	HTTP struct {
		Port string `env-required:"true" yaml:"port" env:"HTTP_PORT"`
	}

	// PG -.
	PG struct {
		PoolMax int    `env-required:"true" yaml:"pool_max" env:"PG_POOL_MAX"`
		URL     string `env-required:"true"                 env:"PG_URL"`
	}

	// S3
	S3 struct {
		ACCESS_KEY        string `env-required:"true" env:"S3_ACCESS_KEY"`
		SECRET_ACCESS_KEY string `env-required:"true" env:"S3_SECRET_ACCESS_KEY"`
		BUCKET_NAME       string `env-required:"true" env:"S3_BUCKET_NAME"`
		ENDPOINT          string `env-required:"true" env:"S3_ENDPOINT"`
		DEFAULT_ICON_URL  string `env-required:"true" env:"DEFAULT_ICON_URL"`
	}
)

// NewConfig returns app config.
func NewConfig() (*Config, error) {
	cfg := &Config{}

	if os.Getenv("GIN_MODE") != "release" {
		if err := godotenv.Load(".env"); err != nil {
			return nil, errors.New("no .env file found")
		}
	}

	err := cleanenv.ReadConfig("./config/config.yml", cfg)
	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
