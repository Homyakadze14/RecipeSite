package usecases

import (
	"context"
	"errors"
	"net/http"

	"github.com/Homyakadze14/RecipeSite/internal/entities"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	sessionKey string = "session_key"
)

var (
	ErrUnauth = errors.New("not authorize")
)

type sessionStorage interface {
	Create(ctx context.Context, sessionID string, userID int) error
	GetUserID(ctx context.Context, sessionID string) (int, error)
	DeleteByID(ctx context.Context, sessionID string) error
	DeleteByUserID(ctx context.Context, userID int) error
}

type SessionUseCase struct {
	storage sessionStorage
}

func NewSessionUseCase(st sessionStorage) *SessionUseCase {
	return &SessionUseCase{
		storage: st,
	}
}

func (u *SessionUseCase) SessionFromContext(ctx *gin.Context) (*entities.Session, error) {
	sess, ok := ctx.Value(sessionKey).(*entities.Session)
	if !ok {
		return nil, ErrUnauth
	}
	return sess, nil
}

func (u *SessionUseCase) Create(ctx context.Context, userID int) (*entities.Session, error) {
	sui := uuid.New().String()

	err := u.storage.Create(ctx, sui, userID)
	if err != nil {
		return nil, err
	}

	return &entities.Session{
		ID:     sui,
		UserID: userID,
	}, nil
}

func (u *SessionUseCase) GetSession(r *http.Request) (*entities.Session, error) {
	sessionCookie, err := r.Cookie("session_id")
	if err != nil {
		return nil, ErrUnauth
	}

	sess := &entities.Session{}
	sess.UserID, err = u.storage.GetUserID(r.Context(), sessionCookie.Value)
	if err != nil {
		return nil, err
	}

	sess.ID = sessionCookie.Value
	return sess, nil
}

func (u *SessionUseCase) DestroySession(ctx *gin.Context) error {
	sess, err := u.SessionFromContext(ctx)
	if err != nil {
		return err
	}

	return u.storage.DeleteByID(ctx, sess.ID)
}

func (u *SessionUseCase) DestroyAllSessions(ctx context.Context, userID int) error {
	return u.storage.DeleteByUserID(ctx, userID)
}
