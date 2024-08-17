package usecases

import (
	"context"
	"errors"
	"fmt"

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

type msgBrokerRepository interface {
	Send(ctx context.Context, message *entities.RecipeCreationMsg) error
}

type SubscribeUseCases struct {
	storage             subscribeStorage
	msgBrokerRepository msgBrokerRepository
}

func NewSubscribeUsecase(st subscribeStorage, msgBrokerRepo msgBrokerRepository) *SubscribeUseCases {
	return &SubscribeUseCases{
		storage:             st,
		msgBrokerRepository: msgBrokerRepo,
	}
}

func (u *SubscribeUseCases) subscribedToYourself(creatorID, ownerID int) bool {
	return creatorID == ownerID
}

func (u *SubscribeUseCases) Subscribe(ctx context.Context, info *entities.SubscribeInfo) error {
	if u.subscribedToYourself(info.CreatorID, info.SubscriberID) {
		return ErrYourselfSubscribe
	}

	_, err := u.storage.GetID(ctx, info)
	if err != nil {
		if errors.Is(err, ErrSubscribeNotFound) {
			err = u.storage.Subscribe(ctx, info)
			if err != nil {
				if errors.Is(err, ErrUserNotFound) {
					return ErrUserNotFound
				}
				return fmt.Errorf("SubscribeUseCases - Subscribe - u.storage.Subscribe: %w", err)
			}
			return nil
		}
		return err
	}

	return ErrSubscribe
}

func (u *SubscribeUseCases) Unsubscribe(ctx context.Context, info *entities.SubscribeInfo) error {
	if u.subscribedToYourself(info.CreatorID, info.SubscriberID) {
		return ErrYourselfUnsubscribe
	}

	_, err := u.storage.GetID(ctx, info)
	if err != nil {
		if errors.Is(err, ErrSubscribeNotFound) {
			return ErrUnsubscribe
		}
		return err
	}

	err = u.storage.Unsubscribe(ctx, info)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return ErrUserNotFound
		}
		return fmt.Errorf("SubscribeUseCases - Unsubscribe - u.storage.Unsubscribe: %w", err)
	}

	return nil
}

func (u *SubscribeUseCases) SendToMsgBroker(ctx context.Context, message *entities.RecipeCreationMsg) error {
	err := u.msgBrokerRepository.Send(ctx, message)
	if err != nil {
		return fmt.Errorf("SubscribeUseCases - SendToMsgBroker - u.msgBrokerRepository.Send: %w", err)
	}

	return nil
}
