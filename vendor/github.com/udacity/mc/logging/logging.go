package logging

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/udacity/mc/config"
)

var loggers []*logrus.Logger

// FileLogger returns a logrus logger that writes to ~/.mc/errs.log.
func FileLogger() *logrus.Logger {
	if len(loggers) == 0 {
		logger := makeLogger()
		loggers = append(loggers, logger)
	}

	return loggers[0]
}

func makeLogger() *logrus.Logger {
	log := logrus.New()
	file, err := os.OpenFile(config.LogsFilepath(), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err == nil {
		log.Out = file
	}

	// Only log the warning severity or above.
	log.Level = logrus.WarnLevel
	return log
}
