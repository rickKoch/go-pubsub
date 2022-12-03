package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/sirupsen/logrus"
)

func logInit() {
	SetFormatter(logrus.StandardLogger())

	logrus.SetLevel(logrus.DebugLevel)
}

func SetFormatter(logger *logrus.Logger) {
	logger.SetFormatter(&logrus.JSONFormatter{
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  "time",
			logrus.FieldKeyLevel: "severity",
			logrus.FieldKeyMsg:   "message",
		},
	})
}

// tail -f development.log | humanlog for better logging
func newDevelopomentLogger() *logrus.Entry {
  log := logrus.StandardLogger()
	log.SetLevel(logrus.DebugLevel)
	file, err := os.OpenFile("./development.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o666)
	if err != nil {
		fmt.Println("unable to log a file")
		os.Exit(1)
	}
	log.SetOutput(file)
	return logrus.NewEntry(log)
}

func newProductionLogger() *logrus.Entry {
	log := logrus.New()
	log.Out = ioutil.Discard
	log.SetLevel(logrus.ErrorLevel)
	return logrus.NewEntry(log)
}
