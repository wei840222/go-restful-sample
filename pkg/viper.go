package pkg

import (
	"log"
	"strings"

	"github.com/spf13/viper"
)

// LoadConfig load config from file and set auto env binding
func LoadConfig() error {
	viper.AddConfigPath("./configs")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("./")
	viper.SetConfigName("config")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Printf("[Viper] not load config from file: %s", err)
		} else {
			return err
		}
	} else {
		log.Printf("[Viper] load config from file: %s", viper.ConfigFileUsed())
	}
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	return nil
}
