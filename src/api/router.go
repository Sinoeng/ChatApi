package api

import (
	"github.com/gin-gonic/gin"

	"primary/api/middleware/jwt"
	"primary/database"
	docs "primary/docs"
	"primary/pubsub"

	swaggerfiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
)

var db *database.ChatApiDB
var toPSCh chan pubsub.Message

func InitRouter(database *database.ChatApiDB, sendCh chan pubsub.Message) *gin.Engine {
	db = database
    toPSCh = sendCh
	var router = gin.Default()

	v1 := router.Group("/v1")
	UnauthUser := v1.Group("/user")
	ProtectedGroup := v1.Group("/protected")
	ProtectedGroup.Use(jwt.JWT())

	ServerGroup := ProtectedGroup.Group("/server")
	UserGroup := ProtectedGroup.Group("/user")
	MessageGroup := ServerGroup.Group("/message")

	v1.GET("/ping", PingHandler)

	AddUserRoutes(UserGroup)
	AddUnauthUserRoutes(UnauthUser)
	AddServerRoutes(ServerGroup)
	AddMessageRoutes(MessageGroup)

	docs.SwaggerInfo.BasePath = "/v1"
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler)) // add swagger
	return router
}
