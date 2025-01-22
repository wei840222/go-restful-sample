package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

func InitViper() error {
	viper.SetConfigName(ConfigFileName)
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("/etc/" + AppName)
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return fmt.Errorf("fatal error config file: %w", err)
		}
	}
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	return nil
}
