package main

import (
	"os"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	log.Print("hello world")
	log.Info().Str("foo", "bar").Msg("Hello world")
	log.Error().Int("status", 404).Stack().Err(errors.New("bad error")).Msg("Bad error")
}
