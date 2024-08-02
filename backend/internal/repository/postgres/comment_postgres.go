package repo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Homyakadze14/RecipeSite/internal/entities"
	"github.com/Homyakadze14/RecipeSite/internal/usecases"
	"github.com/Homyakadze14/RecipeSite/pkg/postgres"
	"github.com/jackc/pgx/v4"
)

type CommentRepo struct {
	*postgres.Postgres
	us *usecases.UserUseCases
}

func NewCommentRepository(pg *postgres.Postgres, us *usecases.UserUseCases) *CommentRepo {
	return &CommentRepo{pg, us}
}

func (r *CommentRepo) Save(ctx context.Context, cm *entities.Comment) error {
	_, err := r.Pool.Exec(ctx, "INSERT INTO comments(user_id, recipe_id, text, created_at, updated_at) VALUES ($1,$2,$3,$4,$5)",
		cm.UserID, cm.RecipeID, cm.Text, time.Now(), time.Now())
	if err != nil {
		return fmt.Errorf("CommentRepo - Save - r.Pool.Exec: %w", err)
	}

	return nil
}

func (r *CommentRepo) Update(ctx context.Context, cm *entities.CommentUpdate) error {
	_, err := r.Pool.Exec(ctx, "UPDATE comments SET text=$1, updated_at=$2 WHERE id=$3", cm.Text, time.Now(), cm.ID)
	if err != nil {
		return fmt.Errorf("CommentRepo - Update - r.Pool.Exec: %w", err)
	}

	return nil
}

func (r *CommentRepo) Delete(ctx context.Context, cm *entities.CommentDelete) error {
	_, err := r.Pool.Exec(ctx, "DELETE FROM comments WHERE id=$1", cm.ID)
	if err != nil {
		return fmt.Errorf("CommentRepo - Delete - r.Pool.Exec: %w", err)
	}

	return nil
}

func (r *CommentRepo) GetByID(ctx context.Context, id int) (*entities.Comment, error) {
	row := r.Pool.QueryRow(ctx, "SELECT * FROM comments WHERE id=$1", id)
	comment := &entities.Comment{}

	err := row.Scan(&comment.ID, &comment.UserID, &comment.RecipeID, &comment.Text, &comment.CreatedAt, &comment.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, usecases.ErrCommentNotFound
		}
		return nil, fmt.Errorf("CommentRepo - GetByID -  row.Scan: %w", err)
	}
	comment.Author, err = r.us.GetAuthor(ctx, comment.UserID)
	if err != nil {
		return nil, fmt.Errorf("CommentRepo - GetByID - r.us.GetAuthor: %w", err)
	}

	return comment, nil
}

func (r *CommentRepo) GetAll(ctx context.Context, recipeID int) ([]entities.Comment, error) {
	rows, err := r.Pool.Query(ctx, "SELECT * FROM comments WHERE recipe_id=$1", recipeID)
	if err != nil {
		return nil, fmt.Errorf("CommentRepo - GetCommets - r.Pool.Query: %w", err)
	}

	comments := make([]entities.Comment, 0, constArraySize)
	for rows.Next() {
		comment := entities.Comment{}
		err = rows.Scan(&comment.ID, &comment.UserID, &comment.RecipeID, &comment.Text, &comment.CreatedAt, &comment.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("CommentRepo - GetCommets - rows.Scan: %w", err)
		}
		comment.Author, err = r.us.GetAuthor(ctx, comment.UserID)
		if err != nil {
			return nil, fmt.Errorf("CommentRepo - GetCommets - r.us.GetAuthor: %w", err)
		}
		comments = append(comments, comment)
	}

	return comments, nil
}
