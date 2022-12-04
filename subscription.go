package main

import (
	"context"
	"regexp"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type Connection struct {
	id      uuid.UUID
	message chan Message
}

func NewSubscription(logger *logrus.Entry, topicName, subscriptionName string) *Subscription {
	return &Subscription{
		name:    subscriptionName,
		topic:   topicName,
		message: make(chan Message, 10),
		logger:  logger,
	}
}

type Subscription struct {
	name        string
	topic       string
	message     chan Message
	connections []*Connection
	messages    []Message
	logger      *logrus.Entry
}

func (sub *Subscription) validate() error {
	regex, err := regexp.Compile("^[a-zA-Z0-9_-]+$")
	if err != nil {
		return err
	}
	valid := regex.MatchString(sub.name)
	if !valid {
		return errors.New("Subscription name is not correct. Please provide correct subscription name ([a-zA-Z0-9_-]).")
	}
	return nil
}

func (sub *Subscription) MessageProcessor(ctx context.Context) {
	for {
		select {
		case msg := <-sub.message:
			sub.messages = append(sub.messages, msg)
			sub.logger.Infof("Subscription '%s' processed message", sub.name)
			for _, connection := range sub.connections {
				connection.message <- msg
				sub.logger.Infof("Connection message processed")
			}
		case <-ctx.Done():
			return
		}
	}
}

func (sub *Subscription) CreateConnection() (*Connection, error) {
	id := uuid.New()
	connection := &Connection{id: id, message: make(chan Message)}
	sub.connections = append(sub.connections, connection)
	sub.logger.Infof("Connection with id '%s' created", id)
	return connection, nil
}

func (sub *Subscription) RemoveConnection(id uuid.UUID) error {
	for i, con := range sub.connections {
		if con.id == id {
			sub.connections = append(sub.connections[:i], sub.connections[i+1:]...)
			sub.logger.Infof("Connection with id '%s' remvoed", id)
			break
		}
	}
	return nil
}

func (sub *Subscription) PullMessages() []Message {
	return sub.messages
}
