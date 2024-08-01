package rabbitmq

import (
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	_defaultConnAttempts = 20
	_defaultConnTimeout  = 5 * time.Second
)

func New(url string) (*amqp.Connection, error) {
	connAttempts := _defaultConnAttempts
	connTimeout := _defaultConnTimeout

	var conn *amqp.Connection
	var err error

	for connAttempts > 0 {
		conn, err = amqp.Dial(url)
		if err == nil {
			break
		}

		log.Printf("Rabbitmq is trying to connect, attempts left: %d", connAttempts)

		time.Sleep(connTimeout)

		connAttempts--
	}

	if err != nil {
		return nil, fmt.Errorf("rabbitmq - NewRabbitmq - connAttempts == 0: %w", err)
	}

	return conn, nil
}
