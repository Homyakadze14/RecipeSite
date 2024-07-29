package usecases

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/Homyakadze14/RecipeSite/internal/entities"
)

var (
	ErrAlreadyLike    = errors.New("this recipe already liked")
	ErrNotLikedYet    = errors.New("this recipe not liked yet")
	ErrRecipeNotExist = errors.New("recipe not exist")
)

type likeStorage interface {
	IsAlreadyLike(ctx context.Context, like *entities.Like) (bool, error)
	LikesCount(ctx context.Context, recipeID int) (int, error)
	Like(ctx context.Context, like *entities.Like) error
	Unlike(ctx context.Context, like *entities.Like) error
	GetLikedRecipies(ctx context.Context, userID int) ([]entities.Recipe, error)
	RecipeExist(ctx context.Context, recipeID int) (bool, error)
}

type sessionManagerForLike interface {
	GetSession(r *http.Request) (*entities.Session, error)
}

type LikeUseCases struct {
	storage        likeStorage
	sessionManager sessionManagerForLike
}

func NewLikeUsecase(st likeStorage, sm sessionManagerForLike) *LikeUseCases {
	return &LikeUseCases{
		storage:        st,
		sessionManager: sm,
	}
}

func (u *LikeUseCases) IsAlreadyLike(ctx context.Context, like *entities.Like) (bool, error) {
	liked, err := u.storage.IsAlreadyLike(ctx, like)
	if err != nil {
		return false, fmt.Errorf("UserUseCase - IsAlreadyLike - u.storage.IsAlreadyLike: %w", err)
	}

	return liked, nil
}

func (u *LikeUseCases) LikesCount(ctx context.Context, recipeID int) (int, error) {
	likesCount, err := u.storage.LikesCount(ctx, recipeID)
	if err != nil {
		return 0, fmt.Errorf("UserUseCase - LikesCount - u.storage.LikesCount: %w", err)
	}

	return likesCount, nil
}

func (u *LikeUseCases) Like(ctx context.Context, r *http.Request, recipeID int) error {
	// Check is recipe exist
	exist, err := u.storage.RecipeExist(ctx, recipeID)
	if err != nil {
		return fmt.Errorf("LikeUseCase - Like - u.storage.RecipeExist: %w", err)
	}

	if !exist {
		return ErrRecipeNotExist
	}

	// Get session
	sess, err := u.sessionManager.GetSession(r)
	if err != nil {
		return fmt.Errorf("LikeUseCase - Like - u.sessionManager.GetSession(r): %w", err)
	}

	// Form like
	like := &entities.Like{
		UserID:   sess.UserID,
		RecipeID: recipeID,
	}

	// Check
	alreadyLike, err := u.storage.IsAlreadyLike(ctx, like)
	if err != nil {
		return fmt.Errorf("LikeUseCase - Like - u.storage.IsAlreadyLike(ctx, like): %w", err)
	}

	if alreadyLike {
		return ErrAlreadyLike
	}

	// Update db
	err = u.storage.Like(ctx, like)
	if err != nil {
		return fmt.Errorf("LikeUseCase - Like - u.storage.Like(ctx, like): %w", err)
	}

	return nil
}

func (u *LikeUseCases) Unlike(ctx context.Context, r *http.Request, recipeID int) error {
	// Check is recipe exist
	exist, err := u.storage.RecipeExist(ctx, recipeID)
	if err != nil {
		return fmt.Errorf("LikeUseCase - Like - u.storage.RecipeExist: %w", err)
	}

	if !exist {
		return ErrRecipeNotExist
	}

	// Get session
	sess, err := u.sessionManager.GetSession(r)
	if err != nil {
		return fmt.Errorf("LikeUseCase - Unlike - u.sessionManager.GetSession(r): %w", err)
	}

	// Form like
	like := &entities.Like{
		UserID:   sess.UserID,
		RecipeID: recipeID,
	}

	// Check
	alreadyLike, err := u.storage.IsAlreadyLike(ctx, like)
	if err != nil {
		return fmt.Errorf("LikeUseCase - Unlike - u.storage.IsAlreadyLike(ctx, like): %w", err)
	}

	if !alreadyLike {
		return ErrNotLikedYet
	}

	// Update db
	err = u.storage.Unlike(ctx, like)
	if err != nil {
		return fmt.Errorf("LikeUseCase - Unlike - u.storage.Like(ctx, like): %w", err)
	}

	return nil
}

func (u *LikeUseCases) GetLikedRecipies(ctx context.Context, userID int) ([]entities.Recipe, error) {
	return u.storage.GetLikedRecipies(ctx, userID)
}
