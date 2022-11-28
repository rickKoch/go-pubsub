package main

import (
	"context"
	"regexp"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type Message struct {
	data  string
	topic string
}

type Topic struct {
	name string
}

func (topic *Topic) validate() error {
  if topic.name == "" {

  }
	regex, err := regexp.Compile("^[a-zA-Z0-9_-]+$")
	if err != nil {
		return err
	}
	valid := regex.MatchString(topic.name)
	if !valid {
		return errors.New("Topic name is not correct. Please provide correct topic name ([a-zA-Z0-9_-]).")
	}
	return nil
}

type Subscription struct {
	name    string
	topic   string
	message chan Message
}

func (subscription *Subscription) validate() error {
	regex, err := regexp.Compile("^[a-zA-Z0-9_-]+$")
	if err != nil {
		return err
	}
	valid := regex.MatchString(subscription.name)
	if !valid {
		return errors.New("Subscription name is not correct. Please provide correct subscription name ([a-zA-Z0-9_-]).")
	}
	return nil
}

type QueueManager struct {
	topics        map[string]Topic
	subscriptions map[string][]*Subscription
	message       chan Message
	logger        *logrus.Entry
}

func NewQueueManager(ctx context.Context, logger *logrus.Entry) QueueManager {
	qm := QueueManager{
		topics:        make(map[string]Topic),
		subscriptions: make(map[string][]*Subscription),
		message:       make(chan Message, 10),
    logger: logger,
	}

  go qm.MessageProcessor(ctx)
	return qm
}

func (qm *QueueManager) CreateTopic(name string) error {
	topic := Topic{name: name}
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

	subscription := &Subscription{
		name:    subscriptionName,
		topic:   topicName,
		message: make(chan Message, 10),
	}
	if err := subscription.validate(); err != nil {
		return err
	}

	qm.subscriptions[topicName] = append(qm.subscriptions[topicName], subscription)
  qm.logger.Infof("Subscription '%s' created", subscriptionName)
	return nil
}

func (qm *QueueManager) Emit(topic string, message string) {
	go func() {
		qm.message <- Message{data: message, topic: topic}
    qm.logger.Infof("Message on '%s' topic emitted", topic)
	}()
}

func (qm *QueueManager) MessageProcessor(ctx context.Context) {
	for {
		select {
		case msg := <-qm.message:
			_, exists := qm.subscriptions[msg.topic]
			if exists {
				for _, sub := range qm.subscriptions[msg.topic] {
					sub.message <- msg
          qm.logger.Infof("Message on '%s' subscription processed", sub.name)
				}
			}
		case <-ctx.Done():
			return
		}
	}
}
