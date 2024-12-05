package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AddChatRoutes(grp *gin.RouterGroup) {

	chat := grp.Group("/chat")

	chat.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, "pong")
	})

}
