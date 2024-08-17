package usecases

import (
	"context"
	"errors"
	"fmt"

	"github.com/Homyakadze14/RecipeSite/internal/entities"
)

var (
	ErrAlreadyLike = errors.New("this recipe already liked")
	ErrNotLikedYet = errors.New("this recipe not liked yet")
)

type likeStorage interface {
	IsAlreadyLike(ctx context.Context, like *entities.Like) (bool, error)
	LikesCount(ctx context.Context, recipeID int) (int, error)
	Like(ctx context.Context, like *entities.Like) error
	Unlike(ctx context.Context, like *entities.Like) error
}

type LikeUseCase struct {
	storage likeStorage
}

func NewLikeUsecase(st likeStorage) *LikeUseCase {
	return &LikeUseCase{
		storage: st,
	}
}

func (u *LikeUseCase) IsAlreadyLike(ctx context.Context, like *entities.Like) (bool, error) {
	liked, err := u.storage.IsAlreadyLike(ctx, like)
	if err != nil {
		return false, fmt.Errorf("LikeUseCase - IsAlreadyLike - u.storage.IsAlreadyLike: %w", err)
	}

	return liked, nil
}

func (u *LikeUseCase) LikesCount(ctx context.Context, recipeID int) (int, error) {
	likesCount, err := u.storage.LikesCount(ctx, recipeID)
	if err != nil {
		return 0, fmt.Errorf("LikeUseCase - LikesCount - u.storage.LikesCount: %w", err)
	}

	return likesCount, nil
}

func (u *LikeUseCase) checkIsNotLikedYet(ctx context.Context, like *entities.Like) error {
	alreadyLike, err := u.IsAlreadyLike(ctx, like)
	if err != nil {
		return fmt.Errorf("LikeUseCase - checkIsNotLikedYet - u.IsAlreadyLike: %w", err)
	}

	if alreadyLike {
		return ErrAlreadyLike
	}

	return nil
}

func (u *LikeUseCase) Like(ctx context.Context, like *entities.Like) error {
	err := u.checkIsNotLikedYet(ctx, like)
	if err != nil {
		return err
	}

	err = u.storage.Like(ctx, like)
	if err != nil {
		return fmt.Errorf("LikeUseCase - Like - u.storage.Like: %w", err)
	}

	return nil
}

func (u *LikeUseCase) checkIsAlreadyLiked(ctx context.Context, like *entities.Like) error {
	alreadyLike, err := u.IsAlreadyLike(ctx, like)
	if err != nil {
		return fmt.Errorf("LikeUseCase - checkIsAlreadyLiked - u.IsAlreadyLike: %w", err)
	}

	if !alreadyLike {
		return ErrNotLikedYet
	}

	return nil
}

func (u *LikeUseCase) Unlike(ctx context.Context, like *entities.Like) error {
	err := u.checkIsAlreadyLiked(ctx, like)
	if err != nil {
		return err
	}

	err = u.storage.Unlike(ctx, like)
	if err != nil {
		return fmt.Errorf("LikeUseCase - Unlike - u.storage.Unlike: %w", err)
	}

	return nil
}
