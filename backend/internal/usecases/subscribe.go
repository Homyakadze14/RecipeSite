package usecases

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/Homyakadze14/RecipeSite/internal/entities"
)

var (
	ErrYourselfSubscribe   = errors.New("you can't subscribe to yourself")
	ErrYourselfUnsubscribe = errors.New("you can't unsubscribe from yourself")
	ErrSubscribeNotFound   = errors.New("subscribe not found")
	ErrSubscribe           = errors.New("you have already subscribed")
	ErrUnsubscribe         = errors.New("you have already unsubscribed")
)

type subscribeStorage interface {
	Subscribe(ctx context.Context, info *entities.SubscribeInfo) error
	Unsubscribe(ctx context.Context, info *entities.SubscribeInfo) error
	GetID(ctx context.Context, info *entities.SubscribeInfo) (int, error)
}

type sessionManagerForSubscribe interface {
	GetSession(r *http.Request) (*entities.Session, error)
}

type SubscribeUseCases struct {
	storage        subscribeStorage
	sessionManager sessionManagerForSubscribe
}

func NewSubscribeUsecase(st subscribeStorage, sm sessionManagerForSubscribe) *SubscribeUseCases {
	return &SubscribeUseCases{
		storage:        st,
		sessionManager: sm,
	}
}

func (u *SubscribeUseCases) Subscribe(ctx context.Context, creator *entities.SubscribeCreator, r *http.Request) error {
	// Get session
	sess, err := u.sessionManager.GetSession(r)
	if err != nil {
		return fmt.Errorf("SubscribeUseCases - Subscribe - u.sessionManager.GetSession: %w", err)
	}

	// Check subscription to yourself
	if creator.ID == sess.UserID {
		return ErrYourselfSubscribe
	}

	// Form info
	info := &entities.SubscribeInfo{
		CreatorID:    creator.ID,
		SubscriberID: sess.UserID,
	}

	// Check if already subscribe
	_, err = u.storage.GetID(ctx, info)
	if err != nil {
		// Not exist
		if errors.Is(err, ErrSubscribeNotFound) {
			// Subscribe
			err = u.storage.Subscribe(ctx, info)
			if err != nil {
				if errors.Is(err, ErrUserNotFound) {
					return ErrUserNotFound
				}
				return fmt.Errorf("SubscribeUseCases - Subscribe - u.storage.Subscribe: %w", err)
			}
		}
		return err
	}

	return ErrSubscribe
}

func (u *SubscribeUseCases) Unsubscribe(ctx context.Context, creator *entities.SubscribeCreator, r *http.Request) error {
	// Get session
	sess, err := u.sessionManager.GetSession(r)
	if err != nil {
		return fmt.Errorf("SubscribeUseCases - Unsubscribe - u.sessionManager.GetSession: %w", err)
	}

	// Check subscription to yourself
	if creator.ID == sess.UserID {
		return ErrYourselfUnsubscribe
	}

	// Form info
	info := &entities.SubscribeInfo{
		CreatorID:    creator.ID,
		SubscriberID: sess.UserID,
	}

	// Check if already subscribe
	_, err = u.storage.GetID(ctx, info)
	if err != nil {
		// Not exist
		if errors.Is(err, ErrSubscribeNotFound) {
			return ErrUnsubscribe
		}
		return err
	}

	// Unsubscribe
	err = u.storage.Unsubscribe(ctx, info)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return ErrUserNotFound
		}
		return fmt.Errorf("SubscribeUseCases - Unsubscribe - u.storage.Unsubscribe: %w", err)
	}

	return nil
}
