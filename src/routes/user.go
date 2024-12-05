package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func new(c *gin.Context) {

}

func verify(c *gin.Context) { // verify jwt

}

func login(c *gin.Context) { // issue jwt

}

func remove(c *gin.Context) {

}

func AddUserRoutes(grp *gin.RouterGroup) {

	user := grp.Group("/user")

	user.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, "pong")
	})

}
