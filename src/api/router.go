package api

import (
	"github.com/gin-gonic/gin"

	"primary/api/middleware/jwt"
	"primary/database"
	docs "primary/docs"

	swaggerfiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
)

var db *database.ChatApiDB

func InitRouter(database *database.ChatApiDB) *gin.Engine {
	db = database
	var router = gin.Default()

	router.POST("/login", LoginHandler)
	router.POST("/newuser", NewHandler)
	router.GET("/ping", PingHandler)

	v1 := router.Group("/v1")
	v1.Use(jwt.JWT())

	AddUserRoutes(v1)
	AddServerRoutes(v1)

	docs.SwaggerInfo.BasePath = "/v1"
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler)) // add swagger
	return router
}
