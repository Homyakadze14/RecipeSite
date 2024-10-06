package usecases

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"

	"github.com/Homyakadze14/RecipeSite/internal/common"
	"github.com/Homyakadze14/RecipeSite/internal/entities"
)

var (
	ErrRecipeNotFound      = errors.New("recipe not found")
	ErrEmptyPhotos         = errors.New("photos must be provided")
	ErrComplexityMustBeInt = errors.New("complexitiy must be integer")
	ErrBadOrderField       = errors.New("bad order field")
)

type recipeStorage interface {
	GetAll(ctx context.Context) ([]entities.Recipe, error)
	GetFiltered(ctx context.Context, filter *entities.RecipeFilter) ([]entities.Recipe, error)
	Get(ctx context.Context, id int) (*entities.Recipe, error)
	Save(ctx context.Context, recipe *entities.Recipe) (id int, err error)
	Update(ctx context.Context, updatedRecipe *entities.Recipe) error
	Delete(ctx context.Context, recipe *entities.Recipe) error
}

type userUseCase interface {
	GetAuthor(ctx context.Context, userID int) (*entities.Author, error)
	GetByLogin(ctx context.Context, login string) (*entities.User, error)
}

type likeUseCase interface {
	LikesCount(ctx context.Context, recipeID int) (int, error)
	IsAlreadyLike(ctx context.Context, like *entities.Like) (bool, error)
}

type subscribeUseCase interface {
	SendToMsgBroker(ctx context.Context, message *entities.RecipeCreationMsg) error
}

type commentUseCase interface {
	GetAll(ctx context.Context, recipeID int) ([]entities.Comment, error)
}

type fileStorageForRecipe interface {
	Save(photos []io.ReadSeeker, contentType string) (string, error)
	Remove(path string) error
}

type cacheRecipeRepository interface {
	Set(ctx context.Context, key string, value interface{}) error
	Get(ctx context.Context, key string, dest interface{}) error
	Del(ctx context.Context, key string) (res int64, err error)
}

type RecipeUseCases struct {
	storage               recipeStorage
	userUseCase           userUseCase
	likeUseCase           likeUseCase
	commentUseCase        commentUseCase
	fileStorage           fileStorageForRecipe
	subscribeUseCase      subscribeUseCase
	cacheRecipeRepository cacheRecipeRepository
}

func NewRecipeUsecase(st recipeStorage, us userUseCase, lu likeUseCase,
	fs fileStorageForRecipe, cu commentUseCase, subu subscribeUseCase, chRep cacheRecipeRepository) *RecipeUseCases {
	return &RecipeUseCases{
		storage:               st,
		userUseCase:           us,
		likeUseCase:           lu,
		fileStorage:           fs,
		commentUseCase:        cu,
		subscribeUseCase:      subu,
		cacheRecipeRepository: chRep,
	}
}

func (r *RecipeUseCases) GetAll(ctx context.Context) ([]entities.RecipeWithAuthor, error) {
	recipes, err := r.storage.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("RecipeUseCase - GetAll - r.storage.GetAll: %w", err)
	}

	rwa := make([]entities.RecipeWithAuthor, 0, 10)
	for _, recipe := range recipes {
		rc := entities.RecipeWithAuthor{
			ID:           recipe.ID,
			UserID:       recipe.UserID,
			Title:        recipe.Title,
			About:        recipe.About,
			Complexitiy:  recipe.Complexitiy,
			NeedTime:     recipe.NeedTime,
			Ingridients:  recipe.Ingridients,
			Instructions: recipe.Instructions,
			PhotosUrls:   recipe.PhotosUrls,
			CreatedAt:    recipe.CreatedAt,
			UpdatedAt:    recipe.UpdatedAt,
		}

		rc.Author, err = r.GetRecipeAuthor(ctx, recipe.UserID)
		if err != nil {
			return nil, fmt.Errorf("RecipeUseCase - GetAll - r.getRecipeAuthor: %w", err)
		}

		rwa = append(rwa, rc)
	}
	return rwa, nil
}

func (r *RecipeUseCases) GetFiltered(ctx context.Context, filter *entities.RecipeFilter) ([]entities.RecipeWithAuthor, error) {
	recipes, err := r.storage.GetFiltered(ctx, filter)
	if err != nil {
		if errors.Is(err, ErrBadOrderField) {
			return nil, ErrBadOrderField
		}
		return nil, fmt.Errorf("RecipeUseCase - GetFiltered - r.storage.GetFiltered: %w", err)
	}

	rwa := make([]entities.RecipeWithAuthor, 0, 10)
	for _, recipe := range recipes {
		rc := entities.RecipeWithAuthor{
			ID:           recipe.ID,
			UserID:       recipe.UserID,
			Title:        recipe.Title,
			About:        recipe.About,
			Complexitiy:  recipe.Complexitiy,
			NeedTime:     recipe.NeedTime,
			Ingridients:  recipe.Ingridients,
			Instructions: recipe.Instructions,
			PhotosUrls:   recipe.PhotosUrls,
			CreatedAt:    recipe.CreatedAt,
			UpdatedAt:    recipe.UpdatedAt,
		}

		rc.Author, err = r.GetRecipeAuthor(ctx, recipe.UserID)
		if err != nil {
			return nil, fmt.Errorf("RecipeUseCase - GetAll - r.getRecipeAuthor: %w", err)
		}

		rwa = append(rwa, rc)
	}
	return rwa, nil
}

func (r *RecipeUseCases) getRecipeFromCache(ctx context.Context, key string) (*entities.Recipe, error) {
	recipe := &entities.Recipe{}
	err := r.cacheRecipeRepository.Get(ctx, key, recipe)
	if err != nil {
		if errors.Is(err, common.ErrCacheKeyNotFound) {
			return nil, common.ErrCacheKeyNotFound
		}
		return nil, err
	}
	return recipe, nil
}

func (r *RecipeUseCases) getRecipeFromStorage(ctx context.Context, id int) (*entities.Recipe, error) {
	recipe, err := r.storage.Get(ctx, id)
	if err != nil {
		if errors.Is(err, ErrRecipeNotFound) {
			return nil, ErrRecipeNotFound
		}
		return nil, fmt.Errorf("RecipeUseCase - getRecipeFromStorage - r.storage.Get: %w", err)
	}
	return recipe, nil
}

func (r *RecipeUseCases) GetRecipeAuthor(ctx context.Context, userID int) (*entities.Author, error) {
	author, err := r.userUseCase.GetAuthor(ctx, userID)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("RecipeUseCase - getRecipeAuthor - r.userUseCase.GetAuthor: %w", err)
	}
	return author, nil
}

func (r *RecipeUseCases) formCacheKey(recipeID int) string {
	return fmt.Sprintf("recipe:%v", recipeID)
}

func (r *RecipeUseCases) Get(ctx context.Context, id, userID int, authorized bool) (*entities.FullRecipe, error) {
	chacheKey := r.formCacheKey(id)
	recipe, err := r.getRecipeFromCache(ctx, chacheKey)

	if err != nil {
		if errors.Is(err, common.ErrCacheKeyNotFound) {
			recipe, err = r.getRecipeFromStorage(ctx, id)
			if err != nil {
				return nil, fmt.Errorf("RecipeUseCase - Get - r.getRecipeFromStorage: %w", err)
			}

			err = r.cacheRecipeRepository.Set(ctx, chacheKey, recipe)
			if err != nil {
				return nil, fmt.Errorf("RecipeUseCase - Get - r.cacheRecipeRepository.Set: %w", err)
			}
		} else {
			return nil, fmt.Errorf("RecipeUseCase - Get - r.getRecipeFromCache: %w", err)
		}
	}

	fullRecipe := entities.FullRecipe{}
	fullRecipe.Recipe = &entities.RecipeWithAuthor{
		ID:           recipe.ID,
		UserID:       recipe.UserID,
		Title:        recipe.Title,
		About:        recipe.About,
		Complexitiy:  recipe.Complexitiy,
		NeedTime:     recipe.NeedTime,
		Ingridients:  recipe.Ingridients,
		Instructions: recipe.Instructions,
		PhotosUrls:   recipe.PhotosUrls,
		CreatedAt:    recipe.CreatedAt,
		UpdatedAt:    recipe.UpdatedAt,
	}

	fullRecipe.Recipe.Author, err = r.GetRecipeAuthor(ctx, fullRecipe.Recipe.UserID)
	if err != nil {
		return nil, fmt.Errorf("RecipeUseCase - Get - r.getRecipeAuthor: %w", err)
	}

	fullRecipe.Comments, err = r.commentUseCase.GetAll(ctx, recipe.ID)
	if err != nil {
		return nil, fmt.Errorf("RecipeUseCase - Get - r.commentUseCase.GetAll: %w", err)
	}

	fullRecipe.LikesCount, err = r.likeUseCase.LikesCount(ctx, recipe.ID)
	if err != nil {
		return nil, fmt.Errorf("RecipeUseCase - Get - r.likeUseCase.LikesCount: %w", err)
	}

	if authorized {
		like := &entities.Like{
			UserID:   userID,
			RecipeID: fullRecipe.Recipe.ID,
		}
		fullRecipe.IsLiked, err = r.likeUseCase.IsAlreadyLike(ctx, like)
		if err != nil {
			return nil, fmt.Errorf("RecipeUseCase - Get - r.likeUseCase.IsAlreadyLike: %w", err)
		}
	}

	return &fullRecipe, nil
}

func (r *RecipeUseCases) getUser(ctx context.Context, login string) (*entities.User, error) {
	user, err := r.userUseCase.GetByLogin(ctx, login)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("RecipeUseCase - getUser - r.userUseCase.GetByLogin: %w", err)
	}
	return user, nil
}

func (r *RecipeUseCases) Create(ctx context.Context, login string, ownerID int, params *entities.CreateRecipe) error {
	user, err := r.getUser(ctx, login)
	if err != nil {
		return fmt.Errorf("RecipeUseCase - Create - r.getUser: %w", err)
	}

	if !common.HavePermisson(ownerID, user.ID) {
		return common.ErrNoPermissions
	}

	if !params.HavePhotos() {
		return ErrEmptyPhotos
	}

	recipe := params.ToRecipe()
	recipe.UserID = user.ID

	recipe.PhotosUrls, err = r.fileStorage.Save(params.Photos, "image/jpeg")
	if err != nil {
		return fmt.Errorf("RecipeUseCase - Create - r.fileStorage.Save: %w", err)
	}

	id, err := r.storage.Save(ctx, recipe)
	if err != nil {
		storageErr := fmt.Errorf("RecipeUseCase - Create - r.storage.Save: %w", err)

		err := r.fileStorage.Remove(recipe.PhotosUrls)
		if err != nil {
			return fmt.Errorf("%w; RecipeUseCase - Create - r.fileStorage.Remove: %w", storageErr, err)
		}

		return storageErr
	}

	message := &entities.RecipeCreationMsg{
		CreatorID: user.ID,
		RecipeID:  id,
	}

	go func() {
		err = r.subscribeUseCase.SendToMsgBroker(ctx, message)
		if err != nil {
			slog.Error(fmt.Sprintf("RecipeUseCase - Create - r.subscribeUseCase.SendToMsgBroker: %s", err.Error()))
		}
	}()

	return nil
}

func (r *RecipeUseCases) Update(ctx context.Context, login string, ownerID, id int, params *entities.UpdateRecipe) error {
	user, err := r.getUser(ctx, login)
	if err != nil {
		return fmt.Errorf("RecipeUseCase - Update - r.getUser: %w", err)
	}

	if !common.HavePermisson(ownerID, user.ID) {
		return common.ErrNoPermissions
	}

	chacheKey := r.formCacheKey(id)
	recipe, err := r.getRecipeFromCache(ctx, chacheKey)
	if err != nil {
		if errors.Is(err, common.ErrCacheKeyNotFound) {
			recipe, err = r.getRecipeFromStorage(ctx, id)
			if err != nil {
				return fmt.Errorf("RecipeUseCase - Update - r.getRecipeFromStorage: %w", err)
			}
		} else {
			return fmt.Errorf("RecipeUseCase - Update - r.getRecipeFromCache: %w", err)
		}
	}
	params.UpdateValues(recipe)

	oldPhotos := ""
	if params.HavePhotos() {
		oldPhotos = recipe.PhotosUrls
		recipe.PhotosUrls, err = r.fileStorage.Save(params.Photos, "image/jpeg")
		if err != nil {
			return fmt.Errorf("RecipeUseCase - Update - r.fileStorage.Save: %w", err)
		}
	}

	err = r.storage.Update(ctx, recipe)
	if err != nil {
		storageErr := fmt.Errorf("RecipeUseCase - Update - r.storage.Update: %w", err)

		err := r.fileStorage.Remove(recipe.PhotosUrls)
		if err != nil {
			return fmt.Errorf("%w; RecipeUseCase - Update - r.fileStorage.Remove: %w", storageErr, err)
		}

		return storageErr
	}

	_, err = r.cacheRecipeRepository.Del(ctx, chacheKey)
	if err != nil {
		return fmt.Errorf("RecipeUseCase - Update - r.cacheRecipeRepository.Del: %w", err)
	}

	if oldPhotos != "" {
		err = r.fileStorage.Remove(oldPhotos)
		if err != nil {
			return fmt.Errorf("RecipeUseCase - Update - r.fileStorage.Remove: %w", err)
		}
	}

	return nil
}

func (r *RecipeUseCases) Delete(ctx context.Context, login string, ownerID int, id int) error {
	user, err := r.getUser(ctx, login)
	if err != nil {
		return fmt.Errorf("RecipeUseCase - Delete - r.getUser: %w", err)
	}

	if !common.HavePermisson(ownerID, user.ID) {
		return common.ErrNoPermissions
	}

	chacheKey := r.formCacheKey(id)
	recipe, err := r.getRecipeFromCache(ctx, chacheKey)
	if err != nil {
		if errors.Is(err, common.ErrCacheKeyNotFound) {
			recipe, err = r.getRecipeFromStorage(ctx, id)
			if err != nil {
				return fmt.Errorf("RecipeUseCase - Delete - r.getRecipeFromStorage: %w", err)
			}
		} else {
			return fmt.Errorf("RecipeUseCase - Delete - r.getRecipeFromCache: %w", err)
		}
	}

	err = r.storage.Delete(ctx, recipe)
	if err != nil {
		return fmt.Errorf("RecipeUseCase - Delete - r.storage.Delete: %w", err)
	}

	_, err = r.cacheRecipeRepository.Del(ctx, chacheKey)
	if err != nil {
		return fmt.Errorf("RecipeUseCase - Delete - r.cacheRecipeRepository.Del: %w", err)
	}

	return nil
}
