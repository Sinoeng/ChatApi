package api

import (
	"github.com/gin-gonic/gin"

	"primary/database"
	docs "primary/docs"

	swaggerfiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
)

var db *database.ChatApiDB

func InitRouter(database *database.ChatApiDB) *gin.Engine {
	db = database
	var router = gin.Default()

	v1 := router.Group("/v1")
	AddUserRoutes(v1)
	AddPingRoutes(v1)

	docs.SwaggerInfo.BasePath = "/v1"
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler)) // add swagger
	return router
}
