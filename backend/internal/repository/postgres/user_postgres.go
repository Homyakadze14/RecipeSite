package repo

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Homyakadze14/RecipeSite/internal/entities"
	"github.com/Homyakadze14/RecipeSite/internal/usecases"
	"github.com/Homyakadze14/RecipeSite/pkg/postgres"
	"github.com/jackc/pgx/v4"
)

type UserRepo struct {
	*postgres.Postgres
}

func NewUserRepository(pg *postgres.Postgres) *UserRepo {
	return &UserRepo{pg}
}

func (r *UserRepo) Create(ctx context.Context, user *entities.User) (id int, err error) {
	row := r.Pool.QueryRow(ctx, "INSERT INTO users(email, login, password, icon_url, about, created_at) VALUES ($1,$2,$3,$4,$5,$6) RETURNING id",
		user.Email, user.Login, user.Password, user.IconURL, user.About, time.Now())
	err = row.Scan(&id)
	if err != nil {
		if strings.Contains(err.Error(), "повторяющееся значение ключа") {
			return -1, usecases.ErrUserUnique
		}
		return -1, fmt.Errorf("UserRepo - Create - r.Pool.QueryRow: %w", err)
	}
	return id, nil
}

func (r *UserRepo) GetByLogin(ctx context.Context, login string) (*entities.User, error) {
	row := r.Pool.QueryRow(ctx, "SELECT * FROM users WHERE login=$1", login)
	usr := &entities.User{}
	err := row.Scan(&usr.ID, &usr.Email, &usr.Login, &usr.Password, &usr.About, &usr.IconURL, &usr.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, usecases.ErrUserNotFound
		}
		return nil, fmt.Errorf("UserRepo - GetByLogin - r.Pool.QueryRow: %w", err)
	}
	return usr, nil
}

func (r *UserRepo) GetByEmail(ctx context.Context, email string) (*entities.User, error) {
	row := r.Pool.QueryRow(ctx, "SELECT * FROM users WHERE email=$1", email)
	usr := &entities.User{}
	err := row.Scan(&usr.ID, &usr.Email, &usr.Login, &usr.Password, &usr.About, &usr.IconURL, &usr.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, usecases.ErrUserNotFound
		}
		return nil, fmt.Errorf("UserRepo - GetByEmail - r.Pool.QueryRow: %w", err)
	}
	return usr, nil
}

func (r *UserRepo) GetAuthor(ctx context.Context, id int) (*entities.Author, error) {
	row := r.Pool.QueryRow(ctx, "SELECT login, icon_url FROM users WHERE id=$1", id)
	usr := &entities.Author{}
	err := row.Scan(&usr.Login, &usr.IconURL)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, usecases.ErrUserNotFound
		}
		return nil, fmt.Errorf("UserRepo - GetAuthor - r.Pool.QueryRow: %w", err)
	}
	return usr, nil
}

func (r *UserRepo) Update(ctx context.Context, id int, user *entities.UserUpdate) error {
	_, err := r.Pool.Exec(ctx, "UPDATE users SET email=$1, login=$2, icon_url=$3, about=$4 WHERE id=$5",
		user.Email, user.Login, user.IconURL, user.About, id)

	if err != nil {
		if strings.Contains(err.Error(), "повторяющееся значение ключа") {
			return usecases.ErrUserUnique
		}
		return fmt.Errorf("UserRepo - Update - r.Pool.Exec: %w", err)
	}
	return nil
}

func (r *UserRepo) UpdatePassword(ctx context.Context, id int, user *entities.UserPasswordUpdate) error {
	_, err := r.Pool.Exec(ctx, "UPDATE users SET password=$1 WHERE id=$2",
		user.Password, id)
	if err != nil {
		return fmt.Errorf("UserRepo - UpdatePassword - r.Pool.Exec: %w", err)
	}
	return nil
}

func (r *UserRepo) GetRecipes(ctx context.Context, userID int) ([]entities.Recipe, error) {
	rows, err := r.Pool.Query(ctx, "SELECT * FROM recipes WHERE user_id=$1", userID)
	if err != nil {
		return nil, fmt.Errorf("UserRepo - GetRecipes - r.Pool.Query: %w", err)
	}

	recipes := make([]entities.Recipe, 0, constArraySize)
	for rows.Next() {
		var recipe entities.Recipe
		err := rows.Scan(&recipe.ID, &recipe.UserID, &recipe.Title, &recipe.About,
			&recipe.Complexitiy, &recipe.NeedTime, &recipe.Ingridients,
			&recipe.PhotosUrls, &recipe.CreatedAt, &recipe.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("UserRepo - GetRecipes - rows.Scan: %w", err)
		}
		recipes = append(recipes, recipe)
	}

	return recipes, nil
}
