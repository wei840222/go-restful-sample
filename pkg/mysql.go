package pkg

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"time"

	"github.com/mattn/go-isatty"
	"github.com/spf13/viper"
	"go.uber.org/fx"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/prometheus"
)

// NewMySQLClient return *gorm.DB and setup prometheus metrics
func NewMySQLClient(lc fx.Lifecycle) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@(%s)/%s?charset=utf8mb4&parseTime=true&loc=%s&time_zone=%s",
		viper.GetString("mysql.username"),
		viper.GetString("mysql.password"),
		viper.GetString("mysql.host"),
		viper.GetString("mysql.database"),
		url.QueryEscape("Asia/Taipei"),
		url.QueryEscape("'+8:00'"),
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		Logger: logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
			SlowThreshold: 100 * time.Millisecond,
			LogLevel: func() logger.LogLevel {
				if viper.GetBool("debug") {
					return logger.Info
				}
				return logger.Warn
			}(),
			Colorful: viper.GetBool("debug") || isatty.IsTerminal(os.Stdout.Fd()) || isatty.IsCygwinTerminal(os.Stdout.Fd()),
		}),
	})
	if err != nil {
		return nil, err
	}

	if err := db.Use(prometheus.New(prometheus.Config{})); err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxIdleConns(3)
	sqlDB.SetMaxOpenConns(30)
	sqlDB.SetConnMaxLifetime(60 * time.Second)

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return sqlDB.PingContext(ctx)
		},
		OnStop: func(context.Context) error {
			return sqlDB.Close()
		},
	})
	return db, nil
}
