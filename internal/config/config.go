package config

import (
	"encoding/json"
	"os"
)

const DefaultIconURL = "https://s3.timeweb.cloud/9bca4c82-d59b4da0-b233-47cd-bd6a-d8ab89558c39/defaul_icon.png"

type Config struct {
	Address              string `json:"address"`
	Port                 int    `json:"port"`
	DB_Host              string `json:"db_host"`
	DB_Name              string `json:"db_name"`
	DB_User              string `json:"db_user"`
	DB_Password          string `json:"db_password"`
	S3_ACCESS_KEY        string `json:"s3_access_key"`
	S3_SECRET_ACCESS_KEY string `json:"s3_secret_access_key"`
	S3_BUCKET_NAME       string `json:"s3_bucket_name"`
	S3_ENDPOINT          string `json:"s3_endpoint"`
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
