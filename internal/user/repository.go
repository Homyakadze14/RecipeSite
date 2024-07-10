package user

import (
	"context"
	"database/sql"
	"time"
)

type UserRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (ur *UserRepository) Create(ctx context.Context, usr *User) (id int, err error) {
	row := ur.db.QueryRowContext(ctx, "INSERT INTO users(email, login, password, icon_url, about, created_at) VALUES ($1,$2,$3,$4,$5,$6) RETURNING id",
		usr.Email, usr.Login, usr.Password, usr.Icon_URL, usr.About, time.Now())
	err = row.Scan(&id)
	if err != nil {
		return -1, err
	}
	return id, nil
}

func (ur *UserRepository) GetByLogin(ctx context.Context, login string) (*User, error) {
	row := ur.db.QueryRowContext(ctx, "SELECT * FROM users WHERE login=$1", login)
	usr := &User{}
	err := row.Scan(&usr.ID, &usr.Email, &usr.Login, &usr.Password, &usr.Icon_URL, &usr.About, &usr.Created_at)
	if err != nil {
		return nil, err
	}
	return usr, nil
}

func (ur *UserRepository) GetByEmail(ctx context.Context, email string) (*User, error) {
	row := ur.db.QueryRowContext(ctx, "SELECT * FROM users WHERE email=$1", email)
	usr := &User{}
	err := row.Scan(&usr.ID, &usr.Email, &usr.Login, &usr.Password, &usr.Icon_URL, &usr.About, &usr.Created_at)
	if err != nil {
		return nil, err
	}
	return usr, nil
}
