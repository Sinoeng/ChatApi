package api

import (
	"github.com/gin-gonic/gin"
)

func new(c *gin.Context) {
	println("Request body: ")
	println(c.Request.Body)

	println("query: ")
	println(c.Query("name"))
}

func verify(c *gin.Context) { // verify jwt

}

func login(c *gin.Context) { // issue jwt

}

func remove(c *gin.Context) {

}

func AddUserRoutes(grp *gin.RouterGroup) {

	user := grp.Group("/user")

	user.GET("/new", new)

}
