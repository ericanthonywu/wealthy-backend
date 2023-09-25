package routers

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func RouterConfig(routeEngine *gin.Engine) {
	routeEngine.Use(gin.Logger())
	routeEngine.Use(gin.Recovery())

	corsConfig := cors.DefaultConfig()
	corsConfig.AddAllowHeaders("Authorization")
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowMethods = []string{"PUT", "POST", "GET", "DELETE", "PATCH"}

	routeEngine.Use(cors.New(corsConfig))
}
