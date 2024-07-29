package usecases

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/Homyakadze14/RecipeSite/internal/entities"
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
	Create(ctx context.Context, recipe *entities.Recipe) error
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

type commentUseCase interface {
	GetAll(ctx context.Context, recipeID int) ([]entities.Comment, error)
}

type fileStorageForRecipe interface {
	Save(image multipart.File, contentType string) (string, error)
	Remove(path string) error
}

type RecipeUseCases struct {
	storage        recipeStorage
	userUseCase    userUseCase
	likeUseCase    likeUseCase
	sessionManager sessionManagerForRecipe
	commentUseCase commentUseCase
	fileStorage    fileStorageForRecipe
}

func NewRecipeUsecase(st recipeStorage, us userUseCase, lu likeUseCase, sm sessionManagerForLike, fs fileStorageForRecipe, cu commentUseCase) *RecipeUseCases {
	return &RecipeUseCases{
		storage:        st,
		sessionManager: sm,
		userUseCase:    us,
		likeUseCase:    lu,
		fileStorage:    fs,
		commentUseCase: cu,
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
	fullRecipe := &entities.FullRecipe{}

	// Get recipe
	recipe, err := r.storage.Get(ctx, recipeID)
	if err != nil {
		if errors.Is(err, ErrRecipeNotFound) {
			return nil, ErrRecipeNotFound
		}
		return nil, fmt.Errorf("RecipeUseCase - Get - r.storage.Get: %w", err)
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
			RecipeID: recipe.ID,
		}
		fullRecipe.IsLiked, err = r.likeUseCase.IsAlreadyLike(ctx, like)
		if err != nil {
			return nil, fmt.Errorf("RecipeUseCase - Get - r.likeUseCase.IsAlreadyLike: %w", err)
		}
	}

	return fullRecipe, nil
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
		return ErrUserNoPermisions
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
	err = r.storage.Create(ctx, recipe)
	if err != nil {
		errImage := r.fileStorage.Remove(recipe.PhotosUrls)
		if errImage != nil {
			return fmt.Errorf("RecipeUseCase - Create - r.storage.Create: %w; RecipeUseCase - Create - r.fileStorage.Remove: %w", err, errImage)
		}

		return fmt.Errorf("RecipeUseCase - Create - r.storage.Create: %w", err)
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
		return ErrUserNoPermisions
	}

	// Get recipe
	dbRecipe, err := r.storage.Get(ctx, recipeID)
	if err != nil {
		if errors.Is(err, ErrRecipeNotFound) {
			return ErrRecipeNotFound
		}
		return fmt.Errorf("RecipeUseCase - Update - r.storage.Get: %w", err)
	}

	// Change values if they was update
	if updatedRecipe.Complexitiy != 0 {
		dbRecipe.Complexitiy = updatedRecipe.Complexitiy
	}
	if updatedRecipe.Title != "" {
		dbRecipe.Title = updatedRecipe.Title
	}
	if updatedRecipe.About != "" {
		dbRecipe.About = updatedRecipe.About
	}
	if updatedRecipe.NeedTime != "" {
		dbRecipe.NeedTime = updatedRecipe.NeedTime
	}
	if updatedRecipe.Ingridients != "" {
		dbRecipe.Ingridients = updatedRecipe.Ingridients
	}

	// Parse photos if exist and save
	multipartFormData, err := gc.MultipartForm()
	if err != nil {
		return fmt.Errorf("RecipeUseCase - Update - gc.MultipartForm(): %w", err)
	}
	files := multipartFormData.File["photos"]
	oldPhotos := ""

	if len(files) != 0 {
		oldPhotos = dbRecipe.PhotosUrls
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
			dbRecipe.PhotosUrls += url + ";"
		}
	}

	// Update recipe in storage
	err = r.storage.Update(ctx, dbRecipe)
	if err != nil {
		errImage := r.fileStorage.Remove(dbRecipe.PhotosUrls)
		if errImage != nil {
			return fmt.Errorf("RecipeUseCase - Update - r.fileStorage.Remove: %w", errImage)
		}

		return fmt.Errorf("RecipeUseCase - Update - r.storage.Update: %w", err)
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
		return ErrUserNoPermisions
	}

	// Get recipe
	recipe, err := r.storage.Get(ctx, id)
	if err != nil {
		if errors.Is(err, ErrRecipeNotFound) {
			return ErrRecipeNotFound
		}
		return fmt.Errorf("RecipeUseCase - Delete - r.storage.Get: %w", err)
	}

	// Delete recipe
	err = r.storage.Delete(ctx, recipe)
	if err != nil {
		return fmt.Errorf("RecipeUseCase - Delete - r.storage.Delete: %w", err)
	}

	return nil
}
