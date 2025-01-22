package config

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

func InitZerolog() {
	l, _ := zerolog.ParseLevel(viper.GetString(ConfigKeyLogLevel))
	if l == zerolog.NoLevel {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	} else {
		zerolog.SetGlobalLevel(l)
	}
	if viper.GetString(ConfigKeyLogFormat) != "json" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339, NoColor: !viper.GetBool(ConfigKeyLogColor)})
	}
}
