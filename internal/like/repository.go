package like

import (
	"context"
	"database/sql"
	"errors"
)

type LikeRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *LikeRepository {
	return &LikeRepository{
		db: db,
	}
}

func (lr *LikeRepository) IsAlreadyLike(ctx context.Context, like *Like) (bool, error) {
	row := lr.db.QueryRowContext(ctx, "SELECT id FROM likes WHERE user_id=$1 AND recipe_id=$2", like.UserID, like.RecipeID)

	var id int
	err := row.Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}

	return true, err
}

func (lr *LikeRepository) LikesCount(ctx context.Context, recipe_id int) (int, error) {
	row := lr.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM likes WHERE recipe_id=$1", recipe_id)

	var count int
	err := row.Scan(&count)
	if err != nil {
		return -1, err
	}

	return count, nil
}

func (lr *LikeRepository) Like(ctx context.Context, like *Like) error {
	_, err := lr.db.ExecContext(ctx, "INSERT INTO likes(user_id, recipe_id) VALUES ($1, $2)", like.UserID, like.RecipeID)
	return err
}

func (lr *LikeRepository) Unlike(ctx context.Context, like *Like) error {
	_, err := lr.db.ExecContext(ctx, "DELETE FROM likes WHERE user_id=$1 AND recipe_id=$2", like.UserID, like.RecipeID)
	return err
}
