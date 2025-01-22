package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"go.uber.org/fx"

	"github.com/wei840222/go-restful-sample/config"
)

func NewGinLogger(notLogged ...string) gin.HandlerFunc {
	var skip map[string]struct{}

	if length := len(notLogged); length > 0 {
		skip = make(map[string]struct{}, length)
		for _, p := range notLogged {
			skip[p] = struct{}{}
		}
	}

	return func(c *gin.Context) {
		// other handler can change c.Path so:
		path := c.Request.URL.Path
		now := time.Now()
		c.Next()
		latency := time.Since(now)
		if latency > time.Minute {
			latency = latency.Truncate(time.Second)
		}
		status := c.Writer.Status()
		clientIP := c.ClientIP()
		clientUserAgent := c.Request.UserAgent()
		referer := c.Request.Referer()
		dataLength := c.Writer.Size()
		if dataLength < 0 {
			dataLength = 0
		}

		if _, ok := skip[path]; ok {
			return
		}

		entry := map[string]any{
			"host":          c.Request.Host,
			"status":        status,
			"latency":       latency.String(),
			"clientIP":      clientIP,
			"method":        c.Request.Method,
			"path":          path,
			"referer":       referer,
			"responseBytes": dataLength,
			"userAgent":     clientUserAgent,
		}

		msg := fmt.Sprintf("[GIN] %v | %d | %13v | %s | %s | %-7s %#v",
			now.Format("2006/01/02 - 15:04:05"),
			status,
			latency,
			c.Request.Host,
			clientIP,
			c.Request.Method,
			path,
		)

		if len(c.Errors) > 0 {
			log.Error().Fields(entry).Msg(strings.TrimSpace(c.Errors.ByType(gin.ErrorTypePrivate).String()))
		}
		if status >= http.StatusInternalServerError {
			log.Error().Fields(entry).Msg(msg)
		} else if status >= http.StatusBadRequest {
			log.Warn().Fields(entry).Msg(msg)
		} else {
			log.Info().Fields(entry).Msg(msg)
		}
	}
}

func NewGinEngine(lc fx.Lifecycle) *gin.Engine {
	if viper.GetString(config.ConfigKeyGinMode) == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	e := gin.New()
	e.ContextWithFallback = true

	e.Use(NewGinLogger(), gin.Recovery())

	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", viper.GetString(config.ConfigKeyGinHost), viper.GetInt(config.ConfigKeyGinPort)),
		Handler: e,
	}

	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			go func() {
				if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					panic(err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return srv.Shutdown(ctx)
		},
	})

	return e
}
