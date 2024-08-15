package usecases

import (
	"context"
	"errors"
	"fmt"

	"github.com/Homyakadze14/RecipeSite/internal/common"
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

type userUseCaseForComment interface {
	GetAuthor(ctx context.Context, id int) (*entities.Author, error)
}

type CommentUseCase struct {
	storage     commentStorage
	userUseCase userUseCaseForComment
}

func NewCommentUseCase(st commentStorage, us userUseCaseForComment) *CommentUseCase {
	return &CommentUseCase{
		storage:     st,
		userUseCase: us,
	}
}

func (u *CommentUseCase) Save(ctx context.Context, cm *entities.Comment) error {
	err := u.storage.Save(ctx, cm)
	if err != nil {
		return fmt.Errorf("CommentUseCase - Save - u.storage.Save: %w", err)
	}

	return nil
}

func (u *CommentUseCase) Update(ctx context.Context, cm *entities.CommentUpdate, ownerID int) error {
	comment, err := u.storage.GetByID(ctx, cm.ID)
	if err != nil {
		return fmt.Errorf("CommentUseCase - Update - u.storage.GetByID: %w", err)
	}

	if !common.HavePermisson(ownerID, comment.UserID) {
		return ErrNoPermissions
	}

	err = u.storage.Update(ctx, cm)
	if err != nil {
		return fmt.Errorf("CommentUseCase - Update - u.storage.Update: %w", err)
	}

	return nil
}

func (u *CommentUseCase) Delete(ctx context.Context, cm *entities.CommentDelete, ownerID int) error {
	comment, err := u.storage.GetByID(ctx, cm.ID)
	if err != nil {
		return fmt.Errorf("CommentUseCase - Delete - u.storage.GetByID: %w", err)
	}

	if !common.HavePermisson(ownerID, comment.UserID) {
		return ErrNoPermissions
	}

	err = u.storage.Delete(ctx, cm)
	if err != nil {
		return fmt.Errorf("CommentUseCase - Delete - u.storage.Delete: %w", err)
	}

	return nil
}

func (u *CommentUseCase) GetAll(ctx context.Context, recipeID int) ([]entities.Comment, error) {
	comments, err := u.storage.GetAll(ctx, recipeID)
	if err != nil {
		return nil, fmt.Errorf("CommentUseCase - GetAll - u.storage.GetAll: %w", err)
	}

	for _, comment := range comments {
		comment.Author, err = u.userUseCase.GetAuthor(ctx, comment.UserID)
		if err != nil {
			return nil, fmt.Errorf("CommentUseCase - GetAll - u.userUseCase.GetAuthor: %w", err)
		}
	}

	return comments, nil
}
