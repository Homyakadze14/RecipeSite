package usecases

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/Homyakadze14/RecipeSite/internal/entities"
)

var (
	ErrCommentNotFound = errors.New("comment not found")
)

type commentStorage interface {
	Save(ctx context.Context, cm *entities.Comment) error
	Update(ctx context.Context, cm *entities.CommentUpdate) error
	Delete(ctx context.Context, cm *entities.CommentDelete) error
	GetAll(ctx context.Context, recipeID int) ([]entities.Comment, error)
	GetByID(ctx context.Context, id int) (*entities.Comment, error)
}

type CommentUseCases struct {
	storage        commentStorage
	sessionManager sessionManagerForLike
}

func NewCommentUsecase(st commentStorage, sm sessionManagerForLike) *CommentUseCases {
	return &CommentUseCases{
		storage:        st,
		sessionManager: sm,
	}
}

func (u *CommentUseCases) Save(ctx context.Context, cm *entities.Comment) error {
	err := u.storage.Save(ctx, cm)
	if err != nil {
		return fmt.Errorf("CommentUseCases - Save - u.storage.Save: %w", err)
	}

	return nil
}

func (u *CommentUseCases) Update(ctx context.Context, r *http.Request, cm *entities.CommentUpdate) error {
	// Get session
	sess, err := u.sessionManager.GetSession(r)
	if err != nil {
		return fmt.Errorf("RecipeUseCase - Update - r.sessionManager.GetSession: %w", err)
	}

	// Get db comment
	comment, err := u.storage.GetByID(ctx, cm.ID)
	if err != nil {
		return fmt.Errorf("RecipeUseCase - Update - u.storage.GetByID: %w", err)
	}

	// Check who update comment
	if sess.UserID != comment.UserID {
		return ErrUserNoPermisions
	}

	// Update
	err = u.storage.Update(ctx, cm)
	if err != nil {
		return fmt.Errorf("CommentUseCases - Update - u.storage.Update: %w", err)
	}

	return nil
}

func (u *CommentUseCases) Delete(ctx context.Context, r *http.Request, cm *entities.CommentDelete) error {
	// Get session
	sess, err := u.sessionManager.GetSession(r)
	if err != nil {
		return fmt.Errorf("RecipeUseCase - Delete - r.sessionManager.GetSession: %w", err)
	}

	// Get db comment
	comment, err := u.storage.GetByID(ctx, cm.ID)
	if err != nil {
		return fmt.Errorf("RecipeUseCase - Delete - u.storage.GetByID: %w", err)
	}

	// Check who delete comment
	if sess.UserID != comment.UserID {
		return ErrUserNoPermisions
	}

	// Delete
	err = u.storage.Delete(ctx, cm)
	if err != nil {
		return fmt.Errorf("CommentUseCases - Delete - u.storage.Delete: %w", err)
	}

	return nil
}

func (u *CommentUseCases) GetAll(ctx context.Context, recipeID int) ([]entities.Comment, error) {
	comments, err := u.storage.GetAll(ctx, recipeID)
	if err != nil {
		return nil, fmt.Errorf("CommentUseCases - GetCommets - u.storage.GetCommets: %w", err)
	}

	return comments, nil
}
