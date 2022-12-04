package main

import (
	"context"
	"reflect"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestMessageProcessor(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	topic := "test-topic"
	logrusEntry = logrus.NewEntry(logrus.StandardLogger())
	sub := NewSubscription(logrusEntry, topic, "test-sub")

	msg := Message{data: "this is test message", topic: topic, ackId: 1}
	sub.message <- msg
	con, err := sub.CreateConnection()
  if err != nil {
    t.Error("Connection can't be created")
  }

	go sub.MessageProcessor(ctx)

  streamedMsg := <- con.message
  defer cancel()

	if !reflect.DeepEqual(sub.messages[0], msg) {
		t.Errorf("MessageProcessor = %v, want %v", msg, sub.messages[0])
	}


	if !reflect.DeepEqual(streamedMsg, msg) {
		t.Errorf("MessageProcessor = %v, want %v", msg, sub.messages[0])
	}
}

func TestConnection(t *testing.T) {
	topic := "test-topic"

	sub := NewSubscription(logrusEntry, topic, "test-sub")
	con, err := sub.CreateConnection()
	if err != nil {
		t.Error("Connection can't be created")
	}

	if len(sub.connections) != 1 && sub.connections[0] != con {
		t.Error("Connection creation failed.")
	}

	if err := sub.RemoveConnection(con.id); err != nil {
		t.Error("Connection can't be removed")
	}
}

func TestAcknowledgeMessage(t *testing.T) {
	logrusEntry = logrus.NewEntry(logrus.StandardLogger())
	sub := &Subscription{
		name:   "test-sub",
		topic:  "test-topic",
		logger: logrusEntry,
		messages: []Message{
			{
				ackId: 1,
				data:  "test-message",
				topic: "test-topic",
			},
		},
	}

	if len(sub.PullMessages()) != 1 {
		t.Error("Subscription should have one message")
	}

	if err := sub.AcknowledgeMessage(1); err != nil {
		t.Error("Message can't be acknowledged")
	}

	if len(sub.PullMessages()) != 0 {
		t.Error("Subscription should not have any messages")
	}
}
