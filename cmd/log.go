package cmd

import (
	"github.com/rs/zerolog"
	"os"
	"time"
)

func NewLogger() zerolog.Logger {
	loggerOutput := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	logger := zerolog.New(loggerOutput).With().Timestamp().Logger()

	return logger
}
