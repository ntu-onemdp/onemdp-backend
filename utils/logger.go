package utils

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

var Logger *zerolog.Logger

func init() {
	// Configure logger
	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}).
		Level(zerolog.TraceLevel).
		With().
		Timestamp().
		Caller().
		Logger()
	Logger = &logger
}
