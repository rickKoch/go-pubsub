package main

import (
	"regexp"

	"github.com/pkg/errors"
)

type Topic struct {
	name          string
	subscriptions []string
}

func (topic *Topic) AddSubscription(name string) {
	topic.subscriptions = append(topic.subscriptions, name)
}

func (topic *Topic) validate() error {
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
