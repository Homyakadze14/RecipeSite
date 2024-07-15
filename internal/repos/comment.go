package repos

import (
	"context"
	"database/sql"
	"time"

	"github.com/Homyakadze14/RecipeSite/internal/models"
)

type CommentRepository struct {
	db *sql.DB
}

func NewCommentRepository(db *sql.DB) *CommentRepository {
	return &CommentRepository{
		db: db,
	}
}

func (r *CommentRepository) Save(ctx context.Context, cm *models.Comment) error {
	_, err := r.db.ExecContext(ctx, "INSERT INTO comments(user_id, recipe_id, text, created_at, updated_at) VALUES ($1,$2,$3,$4,$5)",
		cm.UserID, cm.RecipeID, cm.Text, time.Now(), time.Now())
	return err
}

func (r *CommentRepository) Update(ctx context.Context, cm *models.CommentUpdate) error {
	_, err := r.db.ExecContext(ctx, "UPDATE comments SET text=$1, updated_at=$2 WHERE id=$3", cm.Text, time.Now(), cm.ID)
	return err
}

func (r *CommentRepository) Delete(ctx context.Context, cm *models.CommentDelete) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM comments WHERE id=$1", cm.ID)
	return err
}

func (r *CommentRepository) GetCommets(ctx context.Context, recipeID int, ur *UserRepository) ([]models.Comment, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT * FROM comments WHERE recipe_id=$1", recipeID)
	if err != nil {
		return nil, err
	}

	comments := make([]models.Comment, 0, 10)
	for rows.Next() {
		comment := models.Comment{}
		err = rows.Scan(&comment.ID, &comment.UserID, &comment.RecipeID, &comment.Text, &comment.CreatedAt, &comment.UpdatedAt)
		if err != nil {
			return nil, err
		}
		comment.Author, err = ur.GetAuthor(ctx, comment.UserID)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	return comments, nil
}
