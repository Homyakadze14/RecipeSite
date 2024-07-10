package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	Address     string `json:"address"`
	Port        int    `json:"port"`
	DB_Host     string `json:"db_host"`
	DB_Name     string `json:"db_name"`
	DB_User     string `json:"db_user"`
	DB_Password string `json:"db_password"`
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
