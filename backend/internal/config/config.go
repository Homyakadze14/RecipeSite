package config

import (
	"errors"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var DefaultIconURL = ""

type Config struct {
	Local                bool
	Address              string
	Port                 int
	DB_Host              string
	DB_Name              string
	DB_User              string
	DB_Password          string
	S3_ACCESS_KEY        string
	S3_SECRET_ACCESS_KEY string
	S3_BUCKET_NAME       string
	S3_ENDPOINT          string
}

func Parse() (*Config, error) {
	if os.Getenv("PROD") == "" {
		if err := godotenv.Load("./configs/.env"); err != nil {
			return nil, errors.New("no .env file found")
		}
	}

	DefaultIconURL = os.Getenv("DEFAULT_ICON_URL")
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		return nil, errors.New("port must be integer")
	}

	cfg := &Config{
		Address:              os.Getenv("ADDRESS"),
		Port:                 port,
		DB_Host:              os.Getenv("DATABASE_HOST"),
		DB_Name:              os.Getenv("POSTGRES_DB"),
		DB_User:              os.Getenv("POSTGRES_USER"),
		DB_Password:          os.Getenv("POSTGRES_PASSWORD"),
		S3_ACCESS_KEY:        os.Getenv("S3_ACCESS_KEY"),
		S3_SECRET_ACCESS_KEY: os.Getenv("S3_SECRET_ACCESS_KEY"),
		S3_BUCKET_NAME:       os.Getenv("S3_BUCKET_NAME"),
		S3_ENDPOINT:          os.Getenv("S3_ENDPOINT"),
	}

	return cfg, nil
}
