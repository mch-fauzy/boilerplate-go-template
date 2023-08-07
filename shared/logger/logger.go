package logger

import (
	"os"
	"time"

	"github.com/evermos/boilerplate-go/configs"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// InitLogger initializes the logger
func InitLogger() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.TraceLevel)

	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	// format the output as needed here

	log.Logger = log.Output(output)
	log.Trace().Msg("Zerolog initialized.")
}

// ErrorWithStack logs and error and its stack trace with custom formatting.
func ErrorWithStack(err error) {
	log.Error().Msgf("%+v", errors.WithStack(err))
}

// SetLogLevel sets the desired log level specified in env var.
func SetLogLevel(config *configs.Config) {
	level, err := zerolog.ParseLevel(config.Server.LogLevel)
	if err != nil {
		level = zerolog.TraceLevel
		log.Trace().Str("loglevel", level.String()).Msg("Environment has no log level set up, using default.")
	} else {
		log.Trace().Str("loglevel", level.String()).Msg("Desired log level detected.")
	}
	zerolog.SetGlobalLevel(level)
}
