package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewServerHandler(c *gin.Context) {
	name := c.Query("name")
	if err := db.InsertNewServer(name); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create a new server"})
	}
	c.JSON(http.StatusOK, gin.H{"status": "server " + name + " created"})
}

func DeleteServerHandler(c *gin.Context) {
	//srv := c.Param("server")
	/*
		if err := db.serv // get id by name

		if err := db.DeleteServer() // delete by id
	*/
}

func ChangeServerNameHandler(c *gin.Context) {

}

//TODO: get servers by user

func AddServerRoutes(grp *gin.RouterGroup) {
	server := grp.Group("/server")

	server.POST("/", NewServerHandler)
	server.PATCH("/:server", ChangeServerNameHandler)
	server.DELETE("/:server", removeHandler) //TODO make so user only can delete self
}
