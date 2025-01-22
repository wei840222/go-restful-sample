package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/ipfans/fxlogger"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	_ "go.uber.org/automaxprocs"
	"go.uber.org/fx"

	"github.com/wei840222/go-restful-sample/config"
	"github.com/wei840222/go-restful-sample/handler"
	"github.com/wei840222/go-restful-sample/storage"
)

var (
	flagReplacer = strings.NewReplacer(".", "-")
)

var rootCmd = &cobra.Command{
	Use:   "hello",
	Short: "Hello is a hello world program",
	Long:  `Hello is a hello world program, it will print hello world`,
	PreRunE: func(cmd *cobra.Command, _ []string) error {
		if err := config.InitViper(); err != nil {
			return err
		}

		viper.BindPFlag(config.ConfigKeyLogLevel, cmd.Flags().Lookup(flagReplacer.Replace(config.ConfigKeyLogLevel)))
		viper.BindPFlag(config.ConfigKeyLogFormat, cmd.Flags().Lookup(flagReplacer.Replace(config.ConfigKeyLogFormat)))
		viper.BindPFlag(config.ConfigKeyLogColor, cmd.Flags().Lookup(flagReplacer.Replace(config.ConfigKeyLogColor)))

		config.InitZerolog()

		b, err := json.Marshal(viper.AllSettings())
		if err != nil {
			return err
		}
		log.Debug().Str("config", string(b)).Msg("config loaded")

		return nil
	},
	Run: func(*cobra.Command, []string) {
		app := fx.New(
			fx.Provide(
				NewGorm,
				NewGinEngine,
				storage.NewTodoStorage,
			),
			fx.Invoke(
				handler.RegisterTodoHandler,
			),
			fx.WithLogger(fxlogger.WithZerolog(log.Logger)),
		)

		app.Run()
	},
}

func main() {
	rootCmd.PersistentFlags().String(flagReplacer.Replace(config.ConfigKeyLogLevel), "info", "Log level")
	rootCmd.PersistentFlags().String(flagReplacer.Replace(config.ConfigKeyLogFormat), "console", "Log format")
	rootCmd.PersistentFlags().Bool(flagReplacer.Replace(config.ConfigKeyLogColor), true, "Log color")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
