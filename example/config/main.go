package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	ConfigFileName     = "config"
	ConfigKeyLogLevel  = "log.level"
	ConfigKeyLogFormat = "log.format"
	ConfigKeyLogColor  = "log.color"
)

var (
	flagReplacer = strings.NewReplacer(".", "-")
	envReplacer  = strings.NewReplacer(".", "_")
)

func loadConfig() {
	viper.SetConfigName(ConfigFileName)
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			fmt.Println("fatal error config file:", err)
			os.Exit(1)
		}
	}
	viper.SetEnvKeyReplacer(envReplacer)
	viper.AutomaticEnv()
}

var rootCmd = &cobra.Command{
	Use:   "hello",
	Short: "Hello is a hello world program",
	Long:  `Hello is a hello world program, it will print hello world`,
	PreRun: func(cmd *cobra.Command, _ []string) {
		viper.BindPFlag(ConfigKeyLogLevel, cmd.Flags().Lookup(flagReplacer.Replace(ConfigKeyLogLevel)))
		viper.BindPFlag(ConfigKeyLogFormat, cmd.Flags().Lookup(flagReplacer.Replace(ConfigKeyLogFormat)))
		viper.BindPFlag(ConfigKeyLogColor, cmd.Flags().Lookup(flagReplacer.Replace(ConfigKeyLogColor)))
	},
	RunE: func(*cobra.Command, []string) error {
		b, _ := json.Marshal(viper.AllSettings())
		fmt.Println(string(b))
		fmt.Println("Hello, World!")
		return nil
	},
}

func main() {
	cobra.OnInitialize(loadConfig)

	rootCmd.PersistentFlags().String(flagReplacer.Replace(ConfigKeyLogLevel), "info", "Log level")
	rootCmd.PersistentFlags().String(flagReplacer.Replace(ConfigKeyLogFormat), "console", "Log format")
	rootCmd.PersistentFlags().Bool(flagReplacer.Replace(ConfigKeyLogColor), true, "Log color")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
