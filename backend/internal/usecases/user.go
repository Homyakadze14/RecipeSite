package usecases

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/Homyakadze14/RecipeSite/internal/common"
	"github.com/Homyakadze14/RecipeSite/internal/entities"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type userStorage interface {
	Create(ctx context.Context, user *entities.User) (id int, err error)
	GetByLogin(ctx context.Context, login string) (*entities.User, error)
	GetByEmail(ctx context.Context, email string) (*entities.User, error)
	Update(ctx context.Context, id int, user *entities.UserUpdate) error
	UpdatePassword(ctx context.Context, id int, user *entities.UserPasswordUpdate) error
	GetRecipes(ctx context.Context, userID int) ([]entities.Recipe, error)
	GetAuthor(ctx context.Context, id int) (*entities.Author, error)
}

type fileStorage interface {
	Save(photos []io.ReadSeeker, contentType string) (string, error)
	Remove(path string) error
}

type sessionManager interface {
	Create(ctx context.Context, userID int) (*entities.Session, error)
	GetSession(r *http.Request) (*entities.Session, error)
	DestroySession(ctx *gin.Context) error
	DestroyAllSessions(ctx context.Context, userID int) error
}

type jwtUseCase interface {
	GenerateJWT(userID int) (*entities.JWTToken, error)
	GetDataFromJWT(inToken *entities.JWTToken) (*entities.JWTData, error)
}

type cache interface {
	Set(ctx context.Context, key string, value interface{}) error
	Get(ctx context.Context, key string, dest interface{}) error
	Del(ctx context.Context, key string) (res int64, err error)
}

type UserUseCases struct {
	storage        userStorage
	fileStorage    fileStorage
	sessionManager sessionManager
	defaultIconUrl string
	jwtUseCase     jwtUseCase
	cache          cache
}

func NewUserUsecase(st userStorage, sm sessionManager, df string,
	fs fileStorage, jwt jwtUseCase, cache cache) *UserUseCases {
	return &UserUseCases{
		storage:        st,
		sessionManager: sm,
		defaultIconUrl: df,
		fileStorage:    fs,
		jwtUseCase:     jwt,
		cache:          cache,
	}
}

var (
	ErrUserUnique        = errors.New("user with this credentials already exists")
	ErrUserNotFound      = errors.New("user not found")
	ErrUserWrongPassword = errors.New("wrong password")
	ErrUserNotImage      = errors.New("icon must be image")
)

func (u *UserUseCases) GenerateJWT(r *http.Request) (*entities.JWTToken, error) {
	// Get session
	sess, err := u.sessionManager.GetSession(r)
	if err != nil {
		return nil, fmt.Errorf("UserUseCase - GenerateJWT - u.sessionManager.GetSession: %w", err)
	}

	// Generate token
	token, err := u.jwtUseCase.GenerateJWT(sess.UserID)
	if err != nil {
		return nil, fmt.Errorf("UserUseCase - GenerateJWT - u.jwtUseCase.GenerateJWT: %w", err)
	}

	return token, nil
}

func (u *UserUseCases) GetDataFromJWT(token *entities.JWTToken) (*entities.JWTData, error) {
	data, err := u.jwtUseCase.GetDataFromJWT(token)
	if err != nil {
		if errors.Is(err, ErrBadToken) {
			return nil, ErrBadToken
		}
		return nil, fmt.Errorf("UserUseCase - GetDataFromJWT - u.jwtUseCase.GetDataFromJWT: %w", err)
	}

	return data, nil
}

func (u *UserUseCases) GetAuthor(ctx context.Context, id int) (*entities.Author, error) {
	// Get author from cache
	cackeKey := fmt.Sprintf("author:%v", id)
	author := &entities.Author{}
	err := u.cache.Get(ctx, cackeKey, author)
	if err != nil {
		if errors.Is(err, common.ErrCacheKeyNotFound) {
			// Get author from db
			author, err = u.storage.GetAuthor(ctx, id)
			if err != nil {
				return nil, fmt.Errorf("UserUseCase - GetAuthor - u.storage.GetAuthor: %w", err)
			}

			// Save to cache
			err = u.cache.Set(ctx, cackeKey, author)
			if err != nil {
				return nil, fmt.Errorf("UserUseCase - GetAuthor - u.cache.Set: %w", err)
			}
		} else {
			return nil, fmt.Errorf("UserUseCase - GetAuthor -  u.cache.Get: %w", err)
		}
	}

	return author, nil
}

func (u *UserUseCases) GetByLogin(ctx context.Context, login string) (*entities.User, error) {
	user, err := u.storage.GetByLogin(ctx, login)
	if err != nil {
		return nil, fmt.Errorf("UserUseCase - GetByLogin - u.storage.GetByLogin: %w", err)
	}

	return user, nil
}

func (u *UserUseCases) Signup(ctx context.Context, user *entities.User) (*http.Cookie, error) {
	// set default icon
	user.IconURL = u.defaultIconUrl

	// Hash password
	cryptPass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("UserUseCase - Signup - GeneratePassword: %w", err)
	}
	user.Password = string(cryptPass)

	// Save to storage
	id, err := u.storage.Create(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("UserUseCase - Signup - u.storage.Create: %w", err)
	}

	// Create session
	sess, err := u.sessionManager.Create(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("UserUseCase - Signup - u.sessionManager.Create: %w", err)
	}

	// set cookie
	cookie := &http.Cookie{
		Name:    "session_id",
		Value:   sess.ID,
		Expires: time.Now().Add(90 * 60 * time.Hour),
		Path:    "/",
	}

	return cookie, nil
}

func (u *UserUseCases) Signin(ctx context.Context, user *entities.UserLogin) (*http.Cookie, string, error) {
	// Get db user
	var dbUser *entities.User
	var err error
	if user.Login != "" {
		dbUser, err = u.storage.GetByLogin(ctx, user.Login)
		if err != nil {
			return nil, "", fmt.Errorf("UserUseCase - Signin - u.storage.GetByLogin: %w", err)
		}
	} else if user.Email != "" {
		dbUser, err = u.storage.GetByEmail(ctx, user.Email)
		if err != nil {
			return nil, "", fmt.Errorf("UserUseCase - Signin - u.storage.GetByEmail: %w", err)
		}
	}

	// Check passwords
	err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password))
	if err != nil {
		return nil, "", ErrUserWrongPassword
	}

	// Create session
	sess, err := u.sessionManager.Create(ctx, dbUser.ID)
	if err != nil {
		return nil, "", fmt.Errorf("UserUseCase - Signin - u.sessionManager.Create: %w", err)
	}

	// set cookie
	cookie := &http.Cookie{
		Name:    "session_id",
		Value:   sess.ID,
		Expires: time.Now().Add(90 * 60 * time.Hour),
		Path:    "/",
	}

	return cookie, dbUser.Login, nil
}

func (u *UserUseCases) Logout(ctx *gin.Context) (*http.Cookie, error) {
	err := u.sessionManager.DestroySession(ctx)
	if err != nil {
		return nil, fmt.Errorf("UserUseCase - Logout - u.sessionManager.DestroyByUserID: %w", err)
	}

	cookie := &http.Cookie{
		Name:    "session_id",
		Value:   "sess.ID",
		Expires: time.Now().AddDate(0, 0, -1),
		Path:    "/",
	}

	return cookie, nil
}

func (u *UserUseCases) Update(gc *gin.Context, login string, user *entities.UserUpdate, r *http.Request) (string, error) {
	ctx := gc.Request.Context()
	// Get user from db
	dbUser, err := u.storage.GetByLogin(ctx, login)
	if err != nil {
		return "", fmt.Errorf("UserUseCase - Update - u.storage.GetByLogin: %w", err)
	}

	// Get session
	sess, err := u.sessionManager.GetSession(r)
	if err != nil {
		return "", fmt.Errorf("UserUseCase - Update - u.sessionManager.GetSession: %w", err)
	}

	// Check who update user
	if sess.UserID != dbUser.ID {
		return "", common.ErrNoPermissions
	}

	// Icon
	user.IconURL = dbUser.IconURL
	fileHeader, err := gc.FormFile("icon")
	oldIconUrl := ""
	if fileHeader != nil {
		if err != nil {
			return "", fmt.Errorf("UserUseCase - Update - r.FormFile('icon'): %w", err)
		}
		if !strings.Contains(fileHeader.Header.Get("Content-Type"), "image") {
			return "", ErrUserNotImage
		}

		// save file to storage
		file, err := fileHeader.Open()
		if err != nil {
			return "", fmt.Errorf("UserUseCase - Update - file.Open(): %w", err)
		}
		defer file.Close()

		url, err := u.fileStorage.Save(file, "image/jpeg")
		if err != nil {
			return "", fmt.Errorf("UserUseCase - Update - u.fileStorage.Save: %w", err)
		}
		oldIconUrl = dbUser.IconURL
		user.IconURL = url
	}

	// Replace empty values
	if user.Email == "" {
		user.Email = dbUser.Email
	}
	if user.Login == "" {
		user.Login = dbUser.Login
	}

	// Save to storage
	err = u.storage.Update(ctx, dbUser.ID, user)
	if err != nil {
		// Remove new icon
		imgerr := u.fileStorage.Remove(user.IconURL)
		if imgerr != nil {
			return "", fmt.Errorf("UserUseCase - Update - u.storage.Update: %w; UserUseCase - Update - u.fileStorage.Remove(user.Icon_URL): %w", err, imgerr)
		}
		return "", fmt.Errorf("UserUseCase - Update - u.storage.Update: %w", err)
	}

	// Delete author from cache
	_, err = u.cache.Del(ctx, fmt.Sprintf("author:%v", dbUser.ID))
	if err != nil {
		return "", fmt.Errorf("UserUseCase - Update - u.cache.Del: %w", err)
	}

	// Remove old icon if exist
	if oldIconUrl != "" {
		err = u.fileStorage.Remove(oldIconUrl)
		if err != nil {
			return "", fmt.Errorf("UserUseCase - Update - u.fileStorage.Remove(oldIconUrl): %w", err)
		}
	}

	return user.Login, nil
}

func (u *UserUseCases) UpdatePassword(ctx context.Context, login string, user *entities.UserPasswordUpdate, r *http.Request) error {
	// Get user from db
	dbUser, err := u.storage.GetByLogin(ctx, login)
	if err != nil {
		return fmt.Errorf("UserUseCase - UpdatePassword - u.storage.GetByLogin: %w", err)
	}

	// Get session
	sess, err := u.sessionManager.GetSession(r)
	if err != nil {
		return fmt.Errorf("UserUseCase - UpdatePassword - u.sessionManager.GetSession: %w", err)
	}

	// Check who update user
	if sess.UserID != dbUser.ID {
		return common.ErrNoPermissions
	}

	// Hash password
	cryptPass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("UserUseCase - UpdatePassword - GenerateFromPassword: %w", err)
	}
	user.Password = string(cryptPass)

	// Save to storage
	err = u.storage.UpdatePassword(ctx, dbUser.ID, user)
	if err != nil {
		return fmt.Errorf("UserUseCase - UpdatePassword - u.storage.UpdatePassword: %w", err)
	}

	// Destroy all sessions
	err = u.sessionManager.DestroyAllSessions(ctx, dbUser.ID)
	if err != nil {
		return fmt.Errorf("UserUseCase - UpdatePassword - u.sessionManager.DestroyAllSessions: %w", err)
	}

	return nil
}

func (u *UserUseCases) Get(gc *gin.Context, login string) (*entities.UserInfo, error) {
	r := gc.Request
	ctx := gc.Request.Context()
	// Get user from db
	dbUser, err := u.storage.GetByLogin(ctx, login)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("UserUseCase - Get - u.storage.GetByLogin: %w", err)
	}

	// Parse to new structure
	userInfo := &entities.UserInfo{}
	userInfo.ID = dbUser.ID
	userInfo.Login = dbUser.Login
	userInfo.About = dbUser.About
	userInfo.IconURL = dbUser.IconURL
	userInfo.CreatedAt = dbUser.CreatedAt

	// Get session
	sess, err := u.sessionManager.GetSession(r)
	if err == nil {
		// Check who get user
		if sess.UserID == userInfo.ID {
			// Get liked recipes
			//likedRecipes, err := u.likeUseCases.GetLikedRecipies(r.Context(), userInfo.ID)
			if err != nil {
				return nil, fmt.Errorf("UserUseCase - Get - u.likeUseCases.GetLikedRecipies: %w", err)
			}
			//userInfo.LikedRecipies = likedRecipes
		}
	}

	// Get user recipies
	recipies, err := u.storage.GetRecipes(r.Context(), userInfo.ID)
	if err != nil {
		return nil, fmt.Errorf("UserUseCase - Get - u.recipeUsecases.GetAllByUserID: %w", err)
	}
	userInfo.Recipies = recipies

	return userInfo, nil
}
