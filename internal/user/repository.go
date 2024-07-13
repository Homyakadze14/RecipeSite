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
	err := row.Scan(&usr.ID, &usr.Email, &usr.Login, &usr.Password, &usr.About, &usr.Icon_URL, &usr.Created_at)
	if err != nil {
		return nil, err
	}
	return usr, nil
}

func (ur *UserRepository) GetAuthor(ctx context.Context, id int) (*Author, error) {
	row := ur.db.QueryRowContext(ctx, "SELECT login, icon_url FROM users WHERE id=$1", id)
	usr := &Author{}
	err := row.Scan(&usr.Login, &usr.Icon_URL)
	if err != nil {
		return nil, err
	}
	return usr, nil
}

func (ur *UserRepository) GetByEmail(ctx context.Context, email string) (*User, error) {
	row := ur.db.QueryRowContext(ctx, "SELECT * FROM users WHERE email=$1", email)
	usr := &User{}
	err := row.Scan(&usr.ID, &usr.Email, &usr.Login, &usr.Password, &usr.About, &usr.Icon_URL, &usr.Created_at)
	if err != nil {
		return nil, err
	}
	return usr, nil
}

func (ur *UserRepository) Update(ctx context.Context, user_id int, usr *UserUpdate) error {
	_, err := ur.db.ExecContext(ctx, "UPDATE users SET email=$1, login=$2, icon_url=$3, about=$4 WHERE id=$5",
		usr.Email, usr.Login, usr.Icon_URL, usr.About, user_id)
	return err
}

func (ur *UserRepository) UpdatePassword(ctx context.Context, user_id int, usr *UserPasswordUpdate) error {
	_, err := ur.db.ExecContext(ctx, "UPDATE users SET password=$1 WHERE id=$2",
		usr.Password, user_id)
	return err
}
