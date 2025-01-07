package api

import (
	"net/http"
	"primary/api/middleware/authorization"
	"primary/database"
	"primary/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ServerString struct {
	Name string `form:"name" json:"name" xml:"name" binding:"required"`
}

type ServerID struct {
	ID uint `form:"id" json:"id" xml:"id" binding:"required"`
}

func newServerHandler(c *gin.Context) {
	// get the user id of sender
	claims, err := utils.GetClaims(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "jwt error"})
		return
	}
	userid := claims.Userid // all numbers are float64 after demarshal from json

	// extract server name from request
	var serverName ServerString
	if err := c.Bind(&serverName); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "could not get server name from request"})
		return
	}
	name := serverName.Name
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name cannot be empty"})
		return
	}

	// create a new server
	serverID, err := db.InsertNewServer(name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create a new server"})
		return
	}

	// add creator to server as admin
	if err := db.AddUserToServer(userid, serverID, database.ROLE_SERVER_ADMIN); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create a new server"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "server " + name + " created with id " + strconv.FormatUint(uint64(serverID), 10)})
}

func deleteServerHandler(c *gin.Context) {
	//get server id from request
	serverID, err := strconv.ParseUint(c.Param("serverid"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error:": "serverid is not a number"})
		return
	}

	if err := db.DeleteServerByID(uint(serverID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error:": "no such server"})
	}

	c.JSON(http.StatusOK, gin.H{"status": "server successfully deleted"})
}

func changeServerNameHandler(c *gin.Context) {
	serverID, err := strconv.ParseUint(c.Param("serverid"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error:": "serverid is not a number"})
		return
	}

	//get new server name from request
	var newServerName ServerString
	if err := c.Bind(&newServerName); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "could not read server name"})
		return
	}

	if err := db.ChangeServerName(uint(serverID), newServerName.Name); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error:": "no such server"})
	}

	c.JSON(http.StatusOK, gin.H{"status": "server name changed"})
}

func getUsersByServerHandler(c *gin.Context) {
	//get server id from request
	serverID, err := strconv.ParseUint(c.Param("serverid"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error:": "serverid is not a number"})
		return
	}

	users, err := db.GetAllUsersByServerID(uint(serverID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error:": "no such server"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": users})
}

func addUser(c *gin.Context) {
	//get server id from request
	serverID, err := strconv.ParseUint(c.Param("serverid"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error:": "serverid is not a number"})
		return
	}

	var userIDStruct ServerID
	if err := c.Bind(&userIDStruct); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "could not get userid from request"})
		return
	}
	userID := userIDStruct.ID

	if err = db.AddUserToServer(userID, uint(serverID), database.ROLE_NORMAL); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "could not add user to server"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "user added to server"})
}

func AddServerRoutes(grp *gin.RouterGroup) {
	grp.POST("/:serverid/adduser", func(c *gin.Context) {
		authorization.AuthorizeMiddleware(c, db, authorization.CheckGlobalAdmin, authorization.CheckServerAdmin)
	}, addUser)
	grp.GET("/:serverid/users", func(c *gin.Context) {
		authorization.AuthorizeMiddleware(c, db, authorization.CheckGlobalAdmin, authorization.CheckServerMember)
	}, getUsersByServerHandler)
	grp.POST("/new", newServerHandler)
	grp.PATCH("/:serverid/name", func(c *gin.Context) {
		authorization.AuthorizeMiddleware(c, db, authorization.CheckGlobalAdmin, authorization.CheckServerAdmin)
	}, changeServerNameHandler)
	grp.DELETE("/:serverid", func(c *gin.Context) {
		authorization.AuthorizeMiddleware(c, db, authorization.CheckGlobalAdmin, authorization.CheckServerAdmin)
	}, deleteServerHandler)
}
