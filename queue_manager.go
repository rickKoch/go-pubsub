package main

import (
	"context"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type Message struct {
	data  string
	topic string
}

type QueueManager struct {
	topics        map[string]*Topic
	subscriptions map[string]*Subscription
	message       chan Message
	logger        *logrus.Entry
	ctx           context.Context
}

func NewQueueManager(ctx context.Context, logger *logrus.Entry) QueueManager {
	return QueueManager{
		topics:        make(map[string]*Topic),
		subscriptions: make(map[string]*Subscription),
		message:       make(chan Message, 10),
		logger:        logger,
		ctx:           ctx,
	}
}

func (qm *QueueManager) CreateTopic(name string) error {
	topic := &Topic{name: name}
	if err := topic.validate(); err != nil {
		return err
	}

	qm.topics[name] = topic
	qm.logger.Infof("Topic '%s' created", name)
	return nil
}

func (qm *QueueManager) CreateSubscription(topicName, subscriptionName string) error {
	if _, exists := qm.topics[topicName]; !exists {
		return errors.Errorf("Topic with topic name %s, does not exist", topicName)
	}

	subscription := NewSubscription(qm.logger, topicName, subscriptionName)
	if err := subscription.validate(); err != nil {
		return err
	}

	qm.subscriptions[subscriptionName] = subscription
	qm.topics[topicName].addSubscription(subscriptionName)
	qm.logger.Infof("Subscription '%s' created", subscriptionName)

	go qm.subscriptions[subscriptionName].MessageProcessor(qm.ctx)
	qm.logger.Infof("Listening on subscription: '%s'", subscriptionName)
	return nil
}

func (qm *QueueManager) Emit(topic string, message string) {
	go func() {
		qm.message <- Message{data: message, topic: topic}
		qm.logger.Infof("Message on '%s' topic emitted", topic)
	}()
}

func (qm *QueueManager) Run() {
	for {
		select {
		case msg := <-qm.message:
			_, exists := qm.topics[msg.topic]
			if exists {
				for _, subName := range qm.topics[msg.topic].subscriptions {
					sub, exists := qm.subscriptions[subName]
					if exists {
						sub.message <- msg
						qm.logger.Infof("Message on '%s' subscription processed", subName)
					}
				}
			}
		case <-qm.ctx.Done():
			return
		}
	}
}

func (qm *QueueManager) PullMessages(subscriptionName string) ([]Message, error) {
	sub, exists := qm.subscriptions[subscriptionName]
	if !exists {
		return nil, errors.New("Subscription does not exist")
	}

	return sub.messages, nil
}
