package repo

import (
	"context"
	"errors"
	"fmt"

	"github.com/Homyakadze14/RecipeSite/internal/usecases"
	"github.com/Homyakadze14/RecipeSite/pkg/postgres"
	"github.com/jackc/pgx/v4"
)

type SessionRepo struct {
	*postgres.Postgres
}

func NewSessionRepository(pg *postgres.Postgres) *SessionRepo {
	return &SessionRepo{pg}
}

func (r *SessionRepo) Create(ctx context.Context, sessionID string, userID int) error {
	_, err := r.Pool.Exec(ctx, "INSERT INTO sessions (id, user_id) VALUES ($1, $2)", sessionID, userID)
	if err != nil {
		return fmt.Errorf("SessionRepo - Create - r.Pool.Exec: %w", err)
	}

	return nil
}

func (r *SessionRepo) GetUserID(ctx context.Context, sessionID string) (int, error) {
	var userID int
	row := r.Pool.QueryRow(ctx, "SELECT user_id FROM sessions WHERE id = $1", sessionID)
	err := row.Scan(&userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, usecases.ErrUnauth
		}
		return 0, err
	}

	return userID, nil
}

func (r *SessionRepo) DeleteByID(ctx context.Context, sessionID string) error {
	_, err := r.Pool.Exec(ctx, "DELETE FROM sessions WHERE id = $1", sessionID)
	if err != nil {
		return fmt.Errorf("SessionRepo - GetUserID - r.Pool.Exec: %w", err)
	}

	return nil
}

func (r *SessionRepo) DeleteByUserID(ctx context.Context, userID int) error {
	_, err := r.Pool.Exec(ctx, "DELETE FROM sessions WHERE user_id = $1", userID)
	if err != nil {
		return fmt.Errorf("SessionRepo - DeleteByUserID - r.Pool.Exec: %w", err)
	}

	return nil
}
