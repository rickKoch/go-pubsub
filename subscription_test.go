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

	msg := Message{data: "this is test message", topic: topic}
	sub.message <- msg

	go cancel()
	sub.MessageProcessor(ctx)

	if !reflect.DeepEqual(sub.messages[0], msg) {
		t.Errorf("MessageProcessor = %v, want %v", msg, sub.messages[0])
	}
}
