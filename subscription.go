package main

import (
	"context"
	"regexp"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func NewSubscription(logger *logrus.Entry, topicName, subscriptionName string) *Subscription {
	return &Subscription{
		name:    subscriptionName,
		topic:   topicName,
		message: make(chan Message, 10),
    logger: logger,
	}
}

type Subscription struct {
	name     string
	topic    string
	message  chan Message
	messages []Message
	logger   *logrus.Entry
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
			sub.logger.Infof("Subscription message processed")
		case <-ctx.Done():
			return
		}
	}
}
