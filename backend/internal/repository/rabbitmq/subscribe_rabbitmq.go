package rabbitmqrepo

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/Homyakadze14/RecipeSite/internal/entities"
	amqp "github.com/rabbitmq/amqp091-go"
)

type SubscribeRabbitMQRepo struct {
	rmq     *amqp.Connection
	channel *amqp.Channel
	queue   amqp.Queue
}

func NewSubscribeRabbitMQRepository(rabbitmq *amqp.Connection) (*SubscribeRabbitMQRepo, error) {
	ch, que, err := newChannelAndQueue(rabbitmq)

	if err != nil {
		return nil, err
	}

	return &SubscribeRabbitMQRepo{
		rabbitmq,
		ch,
		que,
	}, nil
}

func newChannelAndQueue(rmq *amqp.Connection) (*amqp.Channel, amqp.Queue, error) {
	ch, err := rmq.Channel()
	if err != nil {
		return nil, amqp.Queue{}, fmt.Errorf("SubscribeRabbitMQRepository - newChannelAndQueue - u.rmq.Channel: %w", err)
	}

	q, err := ch.QueueDeclare(
		"new_recipe",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, amqp.Queue{}, fmt.Errorf("SubscribeRabbitMQRepository - newChannelAndQueue - ch.QueueDeclare: %w", err)
	}

	return ch, q, nil
}

func (u *SubscribeRabbitMQRepo) CloseChan() error {
	err := u.channel.Close()
	if err != nil {
		return fmt.Errorf("SubscribeRabbitMQRepository - CloseChan - u.channel.Close: %w", err)
	}
	return nil
}

func (u *SubscribeRabbitMQRepo) Send(ctx context.Context, message *entities.RecipeCreationMsg) error {
	body, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("SubscribeRabbitMQRepository - Send - json.Marshal: %w", err)
	}
	err = u.channel.PublishWithContext(ctx,
		"",
		u.queue.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        []byte(body),
		})
	if err != nil {
		return fmt.Errorf("SubscribeRabbitMQRepository - Send - ch.PublishWithContext: %w", err)
	}

	slog.Info(fmt.Sprintf("Recipe with creator %v and post_id %v has been sent to rmq", message.CreatorID, message.RecipeID))
	return nil
}
