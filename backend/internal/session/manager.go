package session

import (
	"context"
	"database/sql"
	"errors"
	"net/http"

	"github.com/Homyakadze14/RecipeSite/internal/models"
	"github.com/google/uuid"
)

type ctxKey int

const (
	sessionKey ctxKey = 1
)

var (
	ErrUnauth = errors.New("not authorize")
)

type SessionManager struct {
	db *sql.DB
}

func NewSessionManager(db *sql.DB) *SessionManager {
	return &SessionManager{
		db: db,
	}
}

func (sm *SessionManager) SessionFromContext(ctx context.Context) (*models.Session, error) {
	sess, ok := ctx.Value(sessionKey).(*models.Session)
	if !ok {
		return nil, ErrUnauth
	}
	return sess, nil
}

func (sm *SessionManager) Create(ctx context.Context, user_id int) (*models.Session, error) {
	sui := uuid.New().String()

	_, err := sm.db.ExecContext(ctx, "INSERT INTO sessions (id, user_id) VALUES ($1, $2)", sui, user_id)
	if err != nil {
		return nil, err
	}

	return &models.Session{
		ID:     sui,
		UserID: user_id,
	}, nil
}

func (sm *SessionManager) GetSession(r *http.Request) (*models.Session, error) {
	sessionCookie, err := r.Cookie("session_id")
	if err != nil {
		return nil, ErrUnauth
	}

	sess := &models.Session{}
	row := sm.db.QueryRowContext(r.Context(), "SELECT user_id FROM sessions WHERE id = $1", sessionCookie.Value)
	err = row.Scan(&sess.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUnauth
		} else {
			return nil, err
		}
	}

	sess.ID = sessionCookie.Value
	return sess, nil
}

func (sm *SessionManager) DestroySession(ctx context.Context) error {
	sess, err := sm.SessionFromContext(ctx)
	if err != nil {
		return err
	}
	_, err = sm.db.ExecContext(ctx, "DELETE FROM sessions WHERE id = $1", sess.ID)
	if err != nil {
		return err
	}

	return nil
}

func (sm *SessionManager) DestroyAllSessions(ctx context.Context, user_id int) error {
	_, err := sm.db.ExecContext(ctx, "DELETE FROM sessions WHERE user_id = $1", user_id)
	if err != nil {
		return err
	}

	return nil
}
