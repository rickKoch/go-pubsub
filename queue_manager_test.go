package main

import (
	"context"
	"reflect"
	"testing"

	"github.com/sirupsen/logrus"
)

var (
	ctx         = context.Background()
	logrusEntry = logrus.NewEntry(logrus.StandardLogger())
)

func TestCreateQueueManger(t *testing.T) {
	qm := NewQueueManager(ctx, logrusEntry)

	if len(qm.topics) > 0 {
		t.Error("Topics should be empty")
	}
}

func TestCreateTopic(t *testing.T) {
	qm := NewQueueManager(ctx, logrusEntry)

	topicName := "test-topic"

	err := qm.CreateTopic(topicName)
	if err != nil {
		t.Errorf("Topic with %s topic name, failed to be created", topicName)
	}

	_, exists := qm.topics[topicName]
	if !exists {
		t.Errorf("Topic with %s topic name, does not exist", topicName)
	}
}

func TestValidateTopicName(t *testing.T) {
	qm := NewQueueManager(ctx, logrusEntry)

	topicName := "test*topic"

	err := qm.CreateTopic(topicName)

	if err == nil {
		t.Errorf("Should return error if topic name is not correct: %s", topicName)
	}

	if err.Error() != "Topic name is not correct. Please provide correct topic name ([a-zA-Z0-9_-])." {
		t.Errorf("Should provide topic name error message")
	}
}

func TestCreateSubscription(t *testing.T) {
	qm := NewQueueManager(ctx, logrusEntry)

	topicName := "test-topic"
	subName := "test-sub"

	_ = qm.CreateTopic(topicName)

	err := qm.CreateSubscription(topicName, subName)
	if err != nil {
		t.Errorf("Subscription with %s subscription name, failed to be created", subName)
	}

	topic, exists := qm.topics[topicName]
	if !exists || topic.subscriptions[0] != subName {
		t.Errorf("Topic '%s' does not exists or subscription doesn't belong to the topic", subName)
	}

	_, exists = qm.subscriptions[subName]
	if !exists {
		t.Errorf("Subscription with %s subscription name, does not exist", subName)
	}
}

func TestValidateSubscription(t *testing.T) {
	qm := NewQueueManager(ctx, logrusEntry)
	topic := "test-topic"
	sub := "test^sub"

	tests := []struct {
		name string
		exec func() error
		want string
	}{
		{
			name: "Should not create sub if there is no topic",
			want: "Topic with topic name test-topic, does not exist",
			exec: func() error {
				return qm.CreateSubscription(topic, sub)
			},
		},
		{
			name: "Subscription should not be created with invalid name",
			want: "Subscription name is not correct. Please provide correct subscription name ([a-zA-Z0-9_-]).",
			exec: func() error {
				qm.CreateTopic(topic)
				return qm.CreateSubscription(topic, sub)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.exec()

			if !reflect.DeepEqual(err.Error(), tt.want) {
				t.Errorf("TestValidateSubscription = %v, want %v", err.Error(), tt.want)
			}
		})
	}
}

func TestRun(t *testing.T) {
	topicName := "test-topic"
	subName := "test-sub"

	ctx := context.Background()

	qm := NewQueueManager(ctx, logrusEntry)
	qm.CreateTopic(topicName)
	qm.CreateSubscription(topicName, subName)
  sub, err := qm.Subscription(subName)
  if err != nil {
		t.Errorf("Subscription with %s subscription name, failed to be fetched", subName)
  }

	go qm.Run()

	msgTxt := "this is test message"
	qm.Emit(topicName, msgTxt)

	msg := <-sub.message

	if !reflect.DeepEqual(msg.data, msgTxt) {
		t.Errorf("MessageProcessor = %v, want %v", msg.data, msgTxt)
	}
}
