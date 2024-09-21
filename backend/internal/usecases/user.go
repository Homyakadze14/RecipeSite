package usecases

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
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
	Update(ctx context.Context, user *entities.User) error
	UpdatePassword(ctx context.Context, user *entities.User) error
	GetRecipes(ctx context.Context, userID int) ([]entities.Recipe, error)
	GetAuthor(ctx context.Context, id int) (*entities.Author, error)
	GetIconByLogin(ctx context.Context, login string) (*entities.UserIcon, error)
}

type fileStorage interface {
	Save(photos []io.ReadSeeker, contentType string) (string, error)
	Remove(path string) error
}

type sessionManager interface {
	Create(ctx context.Context, userID int) (*entities.Session, error)
	DestroySession(ctx *gin.Context) error
	DestroyAllSessions(ctx context.Context, userID int) error
}

type jwtUseCase interface {
	GenerateJWT(userID int) (*entities.JWTToken, error)
	GetDataFromJWT(inToken *entities.JWTToken) (*entities.JWTData, error)
}

type likeUseCaseForUser interface {
	GetLikedRecipies(ctx context.Context, userID int) ([]entities.Recipe, error)
}

type cache interface {
	Set(ctx context.Context, key string, value interface{}) error
	Get(ctx context.Context, key string, dest interface{}) error
	Del(ctx context.Context, key string) (res int64, err error)
}

type UserUseCase struct {
	storage        userStorage
	fileStorage    fileStorage
	sessionManager sessionManager
	defaultIconUrl string
	jwtUseCase     jwtUseCase
	cache          cache
	likeUseCase    likeUseCaseForUser
}

func NewUserUsecase(st userStorage, sm sessionManager, df string,
	fs fileStorage, jwt jwtUseCase, cache cache, lu likeUseCaseForUser) *UserUseCase {
	return &UserUseCase{
		storage:        st,
		sessionManager: sm,
		defaultIconUrl: df,
		fileStorage:    fs,
		jwtUseCase:     jwt,
		cache:          cache,
		likeUseCase:    lu,
	}
}

var (
	ErrUserUnique        = errors.New("user with this credentials already exists")
	ErrUserNotFound      = errors.New("user not found")
	ErrUserWrongPassword = errors.New("wrong password")
)

func (u *UserUseCase) GenerateJWT(userID int) (*entities.JWTToken, error) {
	token, err := u.jwtUseCase.GenerateJWT(userID)
	if err != nil {
		return nil, fmt.Errorf("UserUseCase - GenerateJWT - u.jwtUseCase.GenerateJWT: %w", err)
	}

	return token, nil
}

func (u *UserUseCase) GetDataFromJWT(token *entities.JWTToken) (*entities.JWTData, error) {
	data, err := u.jwtUseCase.GetDataFromJWT(token)
	if err != nil {
		if errors.Is(err, ErrBadToken) {
			return nil, ErrBadToken
		}
		return nil, fmt.Errorf("UserUseCase - GetDataFromJWT - u.jwtUseCase.GetDataFromJWT: %w", err)
	}

	return data, nil
}

func (u *UserUseCase) formCacheKey(userID int) string {
	return fmt.Sprintf("author:%v", userID)
}

func (u *UserUseCase) getAuthorFromCache(ctx context.Context, key string) (*entities.Author, error) {
	author := &entities.Author{}
	err := u.cache.Get(ctx, key, author)
	if err != nil {
		if errors.Is(err, common.ErrCacheKeyNotFound) {
			return nil, common.ErrCacheKeyNotFound
		}
		return nil, fmt.Errorf("UserUseCase - getAuthorFromCache - r.cache.Get: %w", err)
	}
	return author, nil
}

func (u *UserUseCase) getAuthorFromStorage(ctx context.Context, id int) (*entities.Author, error) {
	author, err := u.storage.GetAuthor(ctx, id)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("UserUseCase - getAuthorFromStorage - r.storage.GetAuthor: %w", err)
	}
	return author, nil
}

func (u *UserUseCase) GetAuthor(ctx context.Context, id int) (*entities.Author, error) {
	cackeKey := u.formCacheKey(id)
	author, err := u.getAuthorFromCache(ctx, cackeKey)
	if err != nil {
		if errors.Is(err, common.ErrCacheKeyNotFound) {
			author, err = u.getAuthorFromStorage(ctx, id)
			if err != nil {
				return nil, fmt.Errorf("UserUseCase - GetAuthor - u.getAuthorFromStorage: %w", err)
			}

			err = u.cache.Set(ctx, cackeKey, author)
			if err != nil {
				return nil, fmt.Errorf("UserUseCase - GetAuthor - u.cache.Set: %w", err)
			}
		} else {
			return nil, fmt.Errorf("UserUseCase - GetAuthor -  u.getAuthorFromCache: %w", err)
		}
	}

	return author, nil
}

func (u *UserUseCase) GetByLogin(ctx context.Context, login string) (*entities.User, error) {
	user, err := u.storage.GetByLogin(ctx, login)
	if err != nil {
		return nil, fmt.Errorf("UserUseCase - GetByLogin - u.storage.GetByLogin: %w", err)
	}

	return user, nil
}

func (u *UserUseCase) hashPassword(password string) (string, error) {
	cryptPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("UserUseCase - hashPassword - bcrypt.GenerateFromPassword: %w", err)
	}
	return string(cryptPass), nil
}

func (u *UserUseCase) Signup(ctx context.Context, user *entities.User) (*http.Cookie, string, error) {
	user.IconURL = u.defaultIconUrl

	var err error
	user.Password, err = u.hashPassword(user.Password)
	if err != nil {
		return nil, "", fmt.Errorf("UserUseCase - Signup - u.hashPassword: %w", err)
	}

	id, err := u.storage.Create(ctx, user)
	if err != nil {
		return nil, "", fmt.Errorf("UserUseCase - Signup - u.storage.Create: %w", err)
	}

	sess, err := u.sessionManager.Create(ctx, id)
	if err != nil {
		return nil, "", fmt.Errorf("UserUseCase - Signup - u.sessionManager.Create: %w", err)
	}

	cookie := &http.Cookie{
		Name:    "session_id",
		Value:   sess.ID,
		Expires: time.Now().Add(90 * 60 * time.Hour),
		Path:    "/",
	}

	return cookie, user.Login, nil
}

func (u *UserUseCase) comparePasswords(first, second string) error {
	err := bcrypt.CompareHashAndPassword([]byte(first), []byte(second))
	if err != nil {
		return ErrUserWrongPassword
	}
	return nil
}

func (u *UserUseCase) Signin(ctx context.Context, params *entities.UserLogin) (*http.Cookie, string, error) {
	var user *entities.User
	var err error
	if params.Login != "" {
		user, err = u.GetByLogin(ctx, params.Login)
		if err != nil {
			return nil, "", fmt.Errorf("UserUseCase - Signin - u.GetByLogin: %w", err)
		}
	} else if params.Email != "" {
		user, err = u.storage.GetByEmail(ctx, params.Email)
		if err != nil {
			return nil, "", fmt.Errorf("UserUseCase - Signin - u.storage.GetByEmail: %w", err)
		}
	}

	err = u.comparePasswords(user.Password, params.Password)
	if err != nil {
		return nil, "", err
	}

	sess, err := u.sessionManager.Create(ctx, user.ID)
	if err != nil {
		return nil, "", fmt.Errorf("UserUseCase - Signin - u.sessionManager.Create: %w", err)
	}

	cookie := &http.Cookie{
		Name:    "session_id",
		Value:   sess.ID,
		Expires: time.Now().Add(90 * 60 * time.Hour),
		Path:    "/",
	}

	return cookie, user.Login, nil
}

func (u *UserUseCase) Logout(ctx *gin.Context) (*http.Cookie, error) {
	err := u.sessionManager.DestroySession(ctx)
	if err != nil {
		return nil, fmt.Errorf("UserUseCase - Logout - u.sessionManager.DestroySession: %w", err)
	}

	cookie := &http.Cookie{
		Name:    "session_id",
		Value:   "sess.ID",
		Expires: time.Now().AddDate(0, 0, -1),
		Path:    "/",
	}

	return cookie, nil
}

func (u *UserUseCase) Update(ctx context.Context, login string, ownerID int, params *entities.UserUpdate) (string, error) {
	user, err := u.GetByLogin(ctx, login)
	if err != nil {
		return "", fmt.Errorf("UserUseCase - Update - u.GetByLogin: %w", err)
	}

	if !common.HavePermisson(ownerID, user.ID) {
		return "", common.ErrNoPermissions
	}

	params.UpdateValues(user)

	oldIconUrl := ""
	if params.Icon != nil {
		url, err := u.fileStorage.Save([]io.ReadSeeker{params.Icon}, "image/jpeg")
		if err != nil {
			return "", fmt.Errorf("UserUseCase - Update - u.fileStorage.Save: %w", err)
		}
		oldIconUrl = user.IconURL
		user.IconURL = url
	}

	err = u.storage.Update(ctx, user)
	if err != nil {
		storageErr := fmt.Errorf("UserUseCase - Update - u.storage.Update: %w", err)

		err := u.fileStorage.Remove(user.IconURL)
		if err != nil {
			return "", fmt.Errorf("%w; UserUseCase - Update - u.fileStorage.Remove: %w", storageErr, err)
		}

		return "", storageErr
	}

	_, err = u.cache.Del(ctx, u.formCacheKey(user.ID))
	if err != nil {
		return "", fmt.Errorf("UserUseCase - Update - u.cache.Del: %w", err)
	}

	if oldIconUrl != "" {
		err = u.fileStorage.Remove(oldIconUrl)
		if err != nil {
			return "", fmt.Errorf("UserUseCase - Update - u.fileStorage.Remove: %w", err)
		}
	}

	return user.Login, nil
}

func (u *UserUseCase) UpdatePassword(ctx context.Context, login string, ownerID int, params *entities.UserPasswordUpdate) error {
	user, err := u.GetByLogin(ctx, login)
	if err != nil {
		return fmt.Errorf("UserUseCase - UpdatePassword - u.GetByLogin: %w", err)
	}

	if !common.HavePermisson(ownerID, user.ID) {
		return common.ErrNoPermissions
	}

	user.Password, err = u.hashPassword(params.Password)
	if err != nil {
		return fmt.Errorf("UserUseCase - UpdatePassword - u.hashPassword: %w", err)
	}

	err = u.storage.UpdatePassword(ctx, user)
	if err != nil {
		return fmt.Errorf("UserUseCase - UpdatePassword - u.storage.UpdatePassword: %w", err)
	}

	err = u.sessionManager.DestroyAllSessions(ctx, user.ID)
	if err != nil {
		return fmt.Errorf("UserUseCase - UpdatePassword - u.sessionManager.DestroyAllSessions: %w", err)
	}

	return nil
}

func (u *UserUseCase) GetIcon(ctx context.Context, login string) (*entities.UserIcon, error) {
	icn, err := u.storage.GetIconByLogin(ctx, login)
	if err != nil {
		return nil, fmt.Errorf("UserUseCase - GetIcon - u.storage.GetIconByLogin: %w", err)
	}

	return icn, nil
}

func (u *UserUseCase) Get(ctx context.Context, login string, ownerID int, authorized bool) (*entities.UserInfo, error) {
	user, err := u.GetByLogin(ctx, login)
	if err != nil {
		return nil, fmt.Errorf("UserUseCase - Get - u.GetByLogin: %w", err)
	}

	userInfo := &entities.UserInfo{}
	userInfo.ID = user.ID
	userInfo.Login = user.Login
	userInfo.About = user.About
	userInfo.IconURL = user.IconURL
	userInfo.CreatedAt = user.CreatedAt

	if authorized {
		if common.HavePermisson(ownerID, user.ID) {
			likedRecipes, err := u.likeUseCase.GetLikedRecipies(ctx, userInfo.ID)
			if err != nil {
				return nil, fmt.Errorf("UserUseCase - Get - u.likeUseCases.GetLikedRecipies: %w", err)
			}
			userInfo.LikedRecipies = likedRecipes
		}
	}

	recipies, err := u.storage.GetRecipes(ctx, userInfo.ID)
	if err != nil {
		return nil, fmt.Errorf("UserUseCase - Get - u.storage.GetRecipes: %w", err)
	}
	userInfo.Recipies = recipies

	return userInfo, nil
}
