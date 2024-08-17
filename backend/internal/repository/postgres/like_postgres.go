package repo

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Homyakadze14/RecipeSite/internal/entities"
	"github.com/Homyakadze14/RecipeSite/internal/usecases"
	"github.com/Homyakadze14/RecipeSite/pkg/postgres"
	"github.com/jackc/pgx/v5"
)

type LikeRepo struct {
	*postgres.Postgres
}

func NewLikeRepository(pg *postgres.Postgres) *LikeRepo {
	return &LikeRepo{pg}
}

func (l *LikeRepo) IsAlreadyLike(ctx context.Context, like *entities.Like) (bool, error) {
	row := l.Pool.QueryRow(ctx, "SELECT id FROM likes WHERE user_id=$1 AND recipe_id=$2", like.UserID, like.RecipeID)

	var id int
	err := row.Scan(&id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("LikeRepo - IsAlreadyLike - r.Pool.QueryRow: %w", err)
	}

	return true, nil
}

func (l *LikeRepo) LikesCount(ctx context.Context, recipeID int) (int, error) {
	row := l.Pool.QueryRow(ctx, "SELECT COUNT(*) FROM likes WHERE recipe_id=$1", recipeID)

	var count int
	err := row.Scan(&count)
	if err != nil {
		return -1, fmt.Errorf("LikeRepo - LikesCount - r.Pool.QueryRow: %w", err)
	}

	return count, nil
}

func (l *LikeRepo) Like(ctx context.Context, like *entities.Like) error {
	_, err := l.Pool.Exec(ctx, "INSERT INTO likes(user_id, recipe_id) VALUES ($1, $2)", like.UserID, like.RecipeID)
	if err != nil {
		if strings.Contains(err.Error(), "ОШИБКА: INSERT или UPDATE в таблице") {
			return usecases.ErrRecipeNotFound
		}
		return fmt.Errorf("LikeRepo - Like - r.Pool.Exec: %w", err)
	}
	return nil
}

func (l *LikeRepo) Unlike(ctx context.Context, like *entities.Like) error {
	_, err := l.Pool.Exec(ctx, "DELETE FROM likes WHERE user_id=$1 AND recipe_id=$2", like.UserID, like.RecipeID)
	if err != nil {
		if strings.Contains(err.Error(), "ОШИБКА: INSERT или UPDATE в таблице") {
			return usecases.ErrRecipeNotFound
		}
		return fmt.Errorf("LikeRepo - Like - r.Pool.Exec: %w", err)
	}
	return nil
}
