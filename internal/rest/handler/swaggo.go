package handler

import (
	"github.com/wei840222/go-restful-sample/api"

	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	swaggoFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Sample API
// @version 0.1
// @description This is a Sample API.

// @contact.name API Support
// @contact.email wei840222@gmail.com

// @tag.name user

//go:generate swag init --generalInfo=swaggo.go --output=../../../api --parseInternal
// RegisterSwaggoHandler register swagger ui handlers
func RegisterSwaggoHandler(engine *gin.Engine) {
	if apiBaseURL, err := url.Parse(viper.GetString("swagger.base.url")); err == nil {
		api.SwaggerInfo.Schemes = []string{apiBaseURL.Scheme}
		api.SwaggerInfo.Host = apiBaseURL.Host
	}
	if viper.GetBool("debug") {
		engine.GET("/api/swagger", func(c *gin.Context) { c.Redirect(http.StatusMovedPermanently, "/api/swagger/index.html") })
		engine.GET("/api/swagger/*any", ginSwagger.WrapHandler(swaggoFiles.Handler))
	}
}
