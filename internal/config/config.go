package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	Address string `json:"address"`
	Port    int    `json:"port"`
}

func Parse() (*Config, error) {
	file, err := os.ReadFile("./configs/config.json")
	if err != nil {
		return nil, err
	}

	var cfg Config
	json.Unmarshal(file, &cfg)
	return &cfg, err
}
