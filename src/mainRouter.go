package src

import (
	"ChatApi/src/routes"

	"github.com/gin-gonic/gin"

	docs "ChatApi/docs"

	swaggerfiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
)

// @BasePath /v1

func InitRouter() *gin.Engine {
	var router = gin.Default()
	docs.SwaggerInfo.BasePath = "/v1"

	v1 := router.Group("/v1")
	routes.AddUserRoutes(v1)
	routes.AddChatRoutes(v1)
	routes.AddPingRoutes(v1)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	return router
}
