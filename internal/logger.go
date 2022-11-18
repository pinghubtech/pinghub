package internal

import (
	"os"

	"github.com/sirupsen/logrus"
)

func NewLogger() *logrus.Logger {
	lvlString := os.Getenv("LOG_LEVEL")
	if lvlString == "" {
		lvlString = "error"
	}

	lvl, err := logrus.ParseLevel(lvlString)
	if err != nil {
		panic(err)
	}

	logger := logrus.New()
	logger.SetLevel(lvl)

	return logger
}
