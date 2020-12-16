package rest

import (
	"github.com/wei840222/go-restful-sample/internal/rest/middleware"

	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

func init() {
	viper.SetDefault("http.server.port", 8080)
}

func NewRouter(lc fx.Lifecycle) *gin.Engine {
	if viper.GetBool("debug") {
		gin.ForceConsoleColor()
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	engine := gin.New()
	engine.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})
	engine.Use(gin.Logger(), middleware.NewPrometheusMetrics(), gin.Recovery())
	engine.GET("/metrics", func(c *gin.Context) {
		promhttp.Handler().ServeHTTP(c.Writer, c.Request)
	})

	serverAddr := fmt.Sprintf(":%d", viper.GetInt("http.server.port"))
	srv := &http.Server{
		Addr:    serverAddr,
		Handler: engine,
	}

	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			go func() {
				if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					panic(err)
				}
			}()
			log.Printf("[GIN] server start and listen (%s)", serverAddr)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Print("[GIN] server graceful stop")
			return srv.Shutdown(ctx)
		},
	})
	return engine
}
