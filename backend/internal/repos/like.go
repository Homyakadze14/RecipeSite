package repos

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Homyakadze14/RecipeSite/internal/models"
)

type LikeRepository struct {
	db *sql.DB
}

func NewLikeRepository(db *sql.DB) *LikeRepository {
	return &LikeRepository{
		db: db,
	}
}

func (lr *LikeRepository) IsAlreadyLike(ctx context.Context, like *models.Like) (bool, error) {
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

func (lr *LikeRepository) Like(ctx context.Context, like *models.Like) error {
	_, err := lr.db.ExecContext(ctx, "INSERT INTO likes(user_id, recipe_id) VALUES ($1, $2)", like.UserID, like.RecipeID)
	return err
}

func (lr *LikeRepository) Unlike(ctx context.Context, like *models.Like) error {
	_, err := lr.db.ExecContext(ctx, "DELETE FROM likes WHERE user_id=$1 AND recipe_id=$2", like.UserID, like.RecipeID)
	return err
}

func (lr *LikeRepository) GetLikedRecipies(ctx context.Context, userID int) ([]models.Recipe, error) {
	rows, err := lr.db.QueryContext(ctx,
		"SELECT recipes.id, title, about, complexitiy, need_time, ingridients, photos_urls, created_at, updated_at FROM likes JOIN recipes ON recipes.id=likes.recipe_id WHERE likes.user_id=$1",
		userID)

	if err != nil {
		return nil, err
	}

	recipes := make([]models.Recipe, 0, 10)
	for rows.Next() {
		var recipe models.Recipe
		err := rows.Scan(&recipe.ID, &recipe.Title, &recipe.About,
			&recipe.Complexitiy, &recipe.NeedTime, &recipe.Ingridients,
			&recipe.PhotosUrls, &recipe.Created_at, &recipe.Updated_at)
		if err != nil {
			return nil, err
		}
		recipes = append(recipes, recipe)
	}

	return recipes, nil
}
