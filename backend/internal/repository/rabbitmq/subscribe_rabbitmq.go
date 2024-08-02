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
	rmq *amqp.Connection
}

func NewSubscribeRabbitMQRepository(rabbitmq *amqp.Connection) *SubscribeRabbitMQRepo {
	return &SubscribeRabbitMQRepo{rabbitmq}
}

func (u *SubscribeRabbitMQRepo) Send(ctx context.Context, message *entities.NewRecipeRMQMessage) error {
	ch, err := u.rmq.Channel()
	if err != nil {
		return fmt.Errorf("SubscribeRabbitMQRepository - Send - u.rmq.Channel: %w", err)
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
		return fmt.Errorf("SubscribeRabbitMQRepository - Send - ch.QueueDeclare: %w", err)
	}

	body, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("SubscribeRabbitMQRepository - Send - json.Marshal: %w", err)
	}
	err = ch.PublishWithContext(ctx,
		"",
		q.Name,
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
