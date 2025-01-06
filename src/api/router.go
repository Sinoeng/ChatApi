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
	v1 := router.Group("/v1")
	v1.POST("/login", LoginHandler)
	v1.POST("/newuser", NewHandler)
	v1.GET("/ping", PingHandler)

	ProtectedGroup := v1.Group("/protected")
	ProtectedGroup.Use(jwt.JWT())

	ServerGroup := ProtectedGroup.Group("/server")
	UserGroup := ProtectedGroup.Group("/user")
	MessageGroup := ServerGroup.Group("/message")

	AddUserRoutes(UserGroup)
	AddServerRoutes(ServerGroup)
	AddMessageRoutes(MessageGroup)

	docs.SwaggerInfo.BasePath = "/v1"
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler)) // add swagger
	return router
}
