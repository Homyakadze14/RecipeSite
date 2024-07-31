package repo

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Homyakadze14/RecipeSite/internal/entities"
	"github.com/Homyakadze14/RecipeSite/internal/usecases"
	"github.com/Homyakadze14/RecipeSite/pkg/postgres"
	"github.com/jackc/pgx/v4"
)

type SubscribeRepo struct {
	*postgres.Postgres
}

func NewSubscribeRepository(pg *postgres.Postgres) *SubscribeRepo {
	return &SubscribeRepo{pg}
}

func (r *SubscribeRepo) Subscribe(ctx context.Context, info *entities.SubscribeInfo) error {
	_, err := r.Pool.Exec(ctx, "INSERT INTO subscriptions(creator_id, subscriber_id) VALUES ($1,$2)", info.CreatorID, info.SubscriberID)
	if err != nil {
		if strings.Contains(err.Error(), "ОШИБКА: INSERT или UPDATE в таблице") {
			return usecases.ErrUserNotFound
		}
		return fmt.Errorf("SubscribeRepo - Subscribe - r.Pool.Exec: %w", err)
	}
	return nil
}

func (r *SubscribeRepo) Unsubscribe(ctx context.Context, info *entities.SubscribeInfo) error {
	_, err := r.Pool.Exec(ctx, "DELETE FROM subscriptions WHERE creator_id=$1 AND subscriber_id=$2", info.CreatorID, info.SubscriberID)
	if err != nil {
		if strings.Contains(err.Error(), "ОШИБКА: INSERT или UPDATE в таблице") {
			return usecases.ErrUserNotFound
		}
		return fmt.Errorf("SubscribeRepo - Unsubscribe - r.Pool.Exec: %w", err)
	}
	return nil
}

func (r *SubscribeRepo) GetID(ctx context.Context, info *entities.SubscribeInfo) (int, error) {
	row := r.Pool.QueryRow(ctx, "SELECT id FROM subscriptions WHERE creator_id=$1 AND subscriber_id=$2", info.CreatorID, info.SubscriberID)

	var id int
	err := row.Scan(&id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, usecases.ErrSubscribeNotFound
		}
		return 0, fmt.Errorf("SubscribeRepo - GetID - r.Pool.QueryRow: %w", err)
	}
	return id, nil
}
