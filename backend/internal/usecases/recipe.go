package usecases

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/Homyakadze14/RecipeSite/internal/entities"
	redisrepo "github.com/Homyakadze14/RecipeSite/internal/repository/redis"
	"github.com/gin-gonic/gin"
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
	Create(ctx context.Context, recipe *entities.Recipe) (id int, err error)
	Update(ctx context.Context, updatedRecipe *entities.Recipe) error
	Delete(ctx context.Context, recipe *entities.Recipe) error
}

type sessionManagerForRecipe interface {
	GetSession(r *http.Request) (*entities.Session, error)
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
	SendToRmq(ctx context.Context, message *entities.NewRecipeRMQMessage) error
}

type commentUseCase interface {
	GetAll(ctx context.Context, recipeID int) ([]entities.Comment, error)
}

type fileStorageForRecipe interface {
	Save(image multipart.File, contentType string) (string, error)
	Remove(path string) error
}

type redisRecipeRepository interface {
	Set(ctx context.Context, key string, value interface{}) error
	Get(ctx context.Context, key string, dest interface{}) error
	Del(ctx context.Context, key string) (res int64, err error)
}

type RecipeUseCases struct {
	storage               recipeStorage
	userUseCase           userUseCase
	likeUseCase           likeUseCase
	sessionManager        sessionManagerForRecipe
	commentUseCase        commentUseCase
	fileStorage           fileStorageForRecipe
	subscribeUseCase      subscribeUseCase
	redisRecipeRepository redisRecipeRepository
}

func NewRecipeUsecase(st recipeStorage, us userUseCase, lu likeUseCase, sm sessionManagerForLike,
	fs fileStorageForRecipe, cu commentUseCase, subu subscribeUseCase, redRep redisRecipeRepository) *RecipeUseCases {
	return &RecipeUseCases{
		storage:               st,
		sessionManager:        sm,
		userUseCase:           us,
		likeUseCase:           lu,
		fileStorage:           fs,
		commentUseCase:        cu,
		subscribeUseCase:      subu,
		redisRecipeRepository: redRep,
	}
}

func (r *RecipeUseCases) GetAll(ctx context.Context) ([]entities.Recipe, error) {
	recipes, err := r.storage.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("RecipeUseCase - GetAll - r.storage.GetAll: %w", err)
	}
	return recipes, nil
}

func (r *RecipeUseCases) GetFiltered(ctx context.Context, filter *entities.RecipeFilter) ([]entities.Recipe, error) {
	recipes, err := r.storage.GetFiltered(ctx, filter)
	if err != nil {
		if errors.Is(err, ErrBadOrderField) {
			return nil, ErrBadOrderField
		}
		return nil, fmt.Errorf("RecipeUseCase - GetFiltered - r.storage.GetFiltered: %w", err)
	}
	return recipes, nil
}

func (r *RecipeUseCases) Get(ctx context.Context, req *http.Request, recipeID int) (*entities.FullRecipe, error) {
	fullRecipe := entities.FullRecipe{}

	// Get recipe from redis
	redisKey := fmt.Sprintf("recipe:%v", recipeID)
	recipe := &entities.Recipe{}
	err := r.redisRecipeRepository.Get(ctx, redisKey, recipe)
	if err != nil {
		if errors.Is(err, redisrepo.ErrRedisKeyNotFound) {
			// Get recipe from db
			recipe, err = r.storage.Get(ctx, recipeID)
			if err != nil {
				if errors.Is(err, ErrRecipeNotFound) {
					return nil, ErrRecipeNotFound
				}
				return nil, fmt.Errorf("RecipeUseCase - Get - r.storage.Get: %w", err)
			}

			// Save to redis
			err = r.redisRecipeRepository.Set(ctx, redisKey, recipe)
			if err != nil {
				return nil, fmt.Errorf("RecipeUseCase - Get - r.redisRecipeRepository.Set: %w", err)
			}
		} else {
			return nil, fmt.Errorf("RecipeUseCase - Get - r.redisRecipeRepository.Get: %w", err)
		}
	}
	fullRecipe.Recipe = recipe

	// Get Author
	author, err := r.userUseCase.GetAuthor(ctx, recipe.UserID)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("RecipeUseCase - Get - r.userUseCase.GetAuthor: %w", err)
	}
	fullRecipe.Author = author

	// Get likes count
	fullRecipe.LikesCount, err = r.likeUseCase.LikesCount(ctx, recipe.ID)
	if err != nil {
		return nil, fmt.Errorf("RecipeUseCase - Get - r.likeUseCase.LikesCount: %w", err)
	}

	// Get comments
	fullRecipe.Comments, err = r.commentUseCase.GetAll(ctx, recipe.ID)
	if err != nil {
		return nil, fmt.Errorf("RecipeUseCase - Get - r.commentUseCase.GetCommets: %w", err)
	}

	// Get session
	sess, err := r.sessionManager.GetSession(req)
	if err == nil {
		// check is user liked this recipe
		like := &entities.Like{
			UserID:   sess.UserID,
			RecipeID: fullRecipe.Recipe.ID,
		}
		fullRecipe.IsLiked, err = r.likeUseCase.IsAlreadyLike(ctx, like)
		if err != nil {
			return nil, fmt.Errorf("RecipeUseCase - Get - r.likeUseCase.IsAlreadyLike: %w", err)
		}
	}

	return &fullRecipe, nil
}

func (r *RecipeUseCases) Create(gc *gin.Context, userLogin string, crRecipe *entities.CreateRecipe) error {
	ctx := gc.Request.Context()
	req := gc.Request

	// Get user from db
	dbUser, err := r.userUseCase.GetByLogin(ctx, userLogin)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return ErrUserNotFound
		}
		return fmt.Errorf("RecipeUseCase - Create - r.userUseCase.GetByLogin: %w", err)
	}

	// Get session
	sess, err := r.sessionManager.GetSession(req)
	if err != nil {
		return fmt.Errorf("RecipeUseCase - Create - r.sessionManager.GetSession: %w", err)
	}

	// Check who create recipe
	if sess.UserID != dbUser.ID {
		return ErrNoPermissions
	}

	recipe := &entities.Recipe{
		UserID:      dbUser.ID,
		Title:       crRecipe.Title,
		About:       crRecipe.About,
		NeedTime:    crRecipe.NeedTime,
		Ingridients: crRecipe.Ingridients,
		Complexitiy: crRecipe.Complexitiy,
	}

	// Parse photos and save
	multipartFormData, err := gc.MultipartForm()
	if err != nil {
		return fmt.Errorf("RecipeUseCase - Create - gc.MultipartForm(): %w", err)
	}
	files := multipartFormData.File["photos"]

	if len(files) == 0 {
		return ErrEmptyPhotos
	}

	for _, fileHeader := range files {
		if !strings.Contains(fileHeader.Header.Get("Content-Type"), "image") {
			return ErrUserNotImage
		}

		// save file to storage
		file, err := fileHeader.Open()
		if err != nil {
			return fmt.Errorf("RecipeUseCase - Create - file.Open(): %w", err)
		}
		defer file.Close()

		url, err := r.fileStorage.Save(file, "image/jpeg")
		if err != nil {
			return fmt.Errorf("RecipeUseCase - Create - u.fileStorage.Save: %w", err)
		}

		recipe.PhotosUrls += url + ";"
	}

	// Save to storage
	id, err := r.storage.Create(ctx, recipe)
	if err != nil {
		errImage := r.fileStorage.Remove(recipe.PhotosUrls)
		if errImage != nil {
			return fmt.Errorf("RecipeUseCase - Create - r.storage.Create: %w; RecipeUseCase - Create - r.fileStorage.Remove: %w", err, errImage)
		}

		return fmt.Errorf("RecipeUseCase - Create - r.storage.Create: %w", err)
	}

	// Send to rmq
	message := &entities.NewRecipeRMQMessage{
		CreatorID: dbUser.ID,
		RecipeID:  id,
	}

	err = r.subscribeUseCase.SendToRmq(ctx, message)
	if err != nil {
		return fmt.Errorf("RecipeUseCase - Create - RMQ - r.subscribeUseCase.SendToRmq: %w", err)
	}

	return nil
}

func (r *RecipeUseCases) Update(gc *gin.Context, userLogin string, recipeID int, updatedRecipe *entities.UpdateRecipe) error {
	ctx := gc.Request.Context()
	req := gc.Request

	// Get user from db
	dbUser, err := r.userUseCase.GetByLogin(ctx, userLogin)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return ErrUserNotFound
		}
		return fmt.Errorf("RecipeUseCase - Update - r.userUseCase.GetByLogin: %w", err)
	}

	// Get session
	sess, err := r.sessionManager.GetSession(req)
	if err != nil {
		return fmt.Errorf("RecipeUseCase - Update - r.sessionManager.GetSession: %w", err)
	}

	// Check who create recipe
	if sess.UserID != dbUser.ID {
		return ErrNoPermissions
	}

	// Get recipe from redis
	redisKey := fmt.Sprintf("recipe:%v", recipeID)
	recipe := &entities.Recipe{}
	err = r.redisRecipeRepository.Get(ctx, redisKey, recipe)
	if err != nil {
		if errors.Is(err, redisrepo.ErrRedisKeyNotFound) {
			// Get recipe from db
			recipe, err = r.storage.Get(ctx, recipeID)
			if err != nil {
				if errors.Is(err, ErrRecipeNotFound) {
					return ErrRecipeNotFound
				}
				return fmt.Errorf("RecipeUseCase - Update - r.storage.Get: %w", err)
			}
		} else {
			return fmt.Errorf("RecipeUseCase - Update - r.redisRecipeRepository.Get: %w", err)
		}
	}

	// Change values if they was update
	if updatedRecipe.Complexitiy != 0 {
		recipe.Complexitiy = updatedRecipe.Complexitiy
	}
	if updatedRecipe.Title != "" {
		recipe.Title = updatedRecipe.Title
	}
	if updatedRecipe.About != "" {
		recipe.About = updatedRecipe.About
	}
	if updatedRecipe.NeedTime != "" {
		recipe.NeedTime = updatedRecipe.NeedTime
	}
	if updatedRecipe.Ingridients != "" {
		recipe.Ingridients = updatedRecipe.Ingridients
	}

	// Parse photos if exist and save
	multipartFormData, err := gc.MultipartForm()
	if err != nil {
		return fmt.Errorf("RecipeUseCase - Update - gc.MultipartForm(): %w", err)
	}
	files := multipartFormData.File["photos"]
	oldPhotos := ""

	if len(files) != 0 {
		oldPhotos = recipe.PhotosUrls
		recipe.PhotosUrls = ""
		for _, fileHeader := range files {
			if !strings.Contains(fileHeader.Header.Get("Content-Type"), "image") {
				return ErrUserNotImage
			}

			// save file to storage
			file, err := fileHeader.Open()
			if err != nil {
				return fmt.Errorf("RecipeUseCase - Update - file.Open(): %w", err)
			}
			defer file.Close()

			url, err := r.fileStorage.Save(file, "image/jpeg")
			if err != nil {
				return fmt.Errorf("RecipeUseCase - Update - u.fileStorage.Save: %w", err)
			}
			recipe.PhotosUrls += url + ";"
		}
	}

	// Update recipe in storage
	err = r.storage.Update(ctx, recipe)
	if err != nil {
		errImage := r.fileStorage.Remove(recipe.PhotosUrls)
		if errImage != nil {
			return fmt.Errorf("RecipeUseCase - Update - r.storage.Update: %w; RecipeUseCase - Update - r.fileStorage.Remove: %w", err, errImage)
		}

		return fmt.Errorf("RecipeUseCase - Update - r.storage.Update: %w", err)
	}

	// Delete recipe from redis
	_, err = r.redisRecipeRepository.Del(ctx, fmt.Sprintf("recipe:%v", recipe.ID))
	if err != nil {
		return fmt.Errorf("RecipeUseCase - Update - r.redisRecipeRepository.Del: %w", err)
	}

	// Delete old photos
	if oldPhotos != "" {
		err = r.fileStorage.Remove(oldPhotos)
		if err != nil {
			return fmt.Errorf("RecipeUseCase - Update - r.fileStorage.Remove: %w", err)
		}
	}

	return nil
}

func (r *RecipeUseCases) Delete(gc *gin.Context, userLogin string, id int) error {
	ctx := gc.Request.Context()
	req := gc.Request

	// Get user from db
	dbUser, err := r.userUseCase.GetByLogin(ctx, userLogin)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return ErrUserNotFound
		}
		return fmt.Errorf("RecipeUseCase - Delete - r.userUseCase.GetByLogin: %w", err)
	}

	// Get session
	sess, err := r.sessionManager.GetSession(req)
	if err != nil {
		return fmt.Errorf("RecipeUseCase - Delete - r.sessionManager.GetSession: %w", err)
	}

	// Check who create recipe
	if sess.UserID != dbUser.ID {
		return ErrNoPermissions
	}

	// Get recipe from redis
	redisKey := fmt.Sprintf("recipe:%v", id)
	recipe := &entities.Recipe{}
	err = r.redisRecipeRepository.Get(ctx, redisKey, recipe)
	if err != nil {
		if errors.Is(err, redisrepo.ErrRedisKeyNotFound) {
			// Get recipe from db
			recipe, err = r.storage.Get(ctx, id)
			if err != nil {
				if errors.Is(err, ErrRecipeNotFound) {
					return ErrRecipeNotFound
				}
				return fmt.Errorf("RecipeUseCase - Delete - r.storage.Get: %w", err)
			}
		} else {
			return fmt.Errorf("RecipeUseCase - Delete - r.redisRecipeRepository.Get: %w", err)
		}
	}

	// Delete recipe
	err = r.storage.Delete(ctx, recipe)
	if err != nil {
		return fmt.Errorf("RecipeUseCase - Delete - r.storage.Delete: %w", err)
	}

	// Delete recipe from redis
	_, err = r.redisRecipeRepository.Del(ctx, fmt.Sprintf("recipe:%v", recipe.ID))
	if err != nil {
		return fmt.Errorf("RecipeUseCase - Delete - r.redisRecipeRepository.Del: %w", err)
	}

	return nil
}
