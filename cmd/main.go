package main

import (
	"github.com/wei840222/go-restful-sample/internal/rest"
	"github.com/wei840222/go-restful-sample/internal/rest/handler"
	"github.com/wei840222/go-restful-sample/internal/store"
	"github.com/wei840222/go-restful-sample/pkg"

	"go.uber.org/fx"
)

func init() {
	if err := pkg.LoadConfig(); err != nil {
		panic(err)
	}
}

func main() {
	fx.New(
		fx.Provide(
			pkg.NewMySQLClient,
			store.NewUserStore,
			rest.NewRouter,
		),
		fx.Invoke(
			handler.RegisterSwaggoHandler,
			handler.RegisterUserHandler,
		),
	).Run()
}
