package database

import (
	"database/sql"
	"fmt"

	"github.com/Homyakadze14/RecipeSite/RecipeSite/internal/config"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func New(cfg *config.Config) (*sql.DB, error) {
	dsn := fmt.Sprintf("host=%s dbname=%s user=%s password=%s sslmode=disable", cfg.DB_Host, cfg.DB_Name, cfg.DB_User, cfg.DB_Password)
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
