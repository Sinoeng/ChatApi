package api

import (
	"net/http"
	"primary/database"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

func NewServerHandler(c *gin.Context) { // add user that created it to server
	tokenAny, exists := c.Get("token")
	if !exists {
		return
	}
	token := tokenAny.(*jwt.Token)
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not read token claims"})
		return
	}
	userid := uint(claims["userid"].(float64))

	name := c.Query("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name cannot be empty"})
	}

	serverID, err := db.InsertNewServer(name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create a new server"})
	}

	if err := db.AddUserToServer(userid, serverID, database.ROLE_SERVER_ADMIN); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create a new server"})
	}
	c.JSON(http.StatusOK, gin.H{"status": "server " + name + " created with id " + strconv.FormatUint(uint64(serverID), 10)})
}

func DeleteServerHandler(c *gin.Context) {
	tokenAny, exists := c.Get("token")
	if !exists {
		return
	}
	token := tokenAny.(*jwt.Token)
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not read token claims"})
		return
	}

	serverID, err := strconv.ParseUint(c.Param("serverid"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error:": err.Error()})
		return
	}

	userid := uint(claims["userid"].(float64))

	userServer, err := db.GetUserServerByIDs(userid, uint(serverID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error:": err.Error()})
		return
	}
	if userServer.Role != database.ROLE_SERVER_ADMIN {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "you do not have permission to delete this server"})
		return
	}

	if err := db.DeleteServerByID(uint(serverID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error:": err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{"status": "server successfully deleted"})
}

func ChangeServerNameHandler(c *gin.Context) {
	tokenAny, exists := c.Get("token")
	if !exists {
		return
	}
	token := tokenAny.(*jwt.Token)
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not read token claims"})
		return
	}
	serverID, err := strconv.ParseUint(c.Param("serverid"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error:": err.Error()})
		return
	}
	userid := uint(claims["userid"].(float64))
	userServer, err := db.GetUserServerByIDs(userid, uint(serverID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error:": err.Error()})
		return
	}
	if userServer.Role != database.ROLE_SERVER_ADMIN {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "you do not have permission to change the name of this server"})
		return
	}
	var serverName ServerName
	if err := c.Bind(&serverName); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := db.ChangeServerName(uint(serverID), serverName.Name); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error:": err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{"status": "server name changed"})
}

type ServerName struct {
	Name string `form:"name" json:"name" xml:"name" binding:"required"`
}

func GetUsersByServerHandler(c *gin.Context) {
	//check that requesting user is memeber or admin
	tokenAny, exists := c.Get("token")
	if !exists {
		return
	}
	token := tokenAny.(*jwt.Token)
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not read token claims"})
		return
	}
	serverID, err := strconv.ParseUint(c.Param("serverid"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error:": err.Error()})
		return
	}
	userid := uint(claims["userid"].(float64))

	_, err = db.GetUserServerByIDs(userid, uint(serverID))
	if err == gorm.ErrRecordNotFound {
		c.JSON(http.StatusBadRequest, gin.H{"error": "not a memeber of this server"})
		return
	}

	server, err := db.GetServerByID(uint(serverID))
	users := server.Users

	c.JSON(http.StatusOK, gin.H{"data": users})
}

func AddServerRoutes(grp *gin.RouterGroup) {
	server := grp.Group("/server")

	server.GET("/users/:serverid", GetUsersByServerHandler)
	server.POST("/new", NewServerHandler)
	server.PATCH("/:serverid", ChangeServerNameHandler)
	server.DELETE("/:serverid", removeHandler) //TODO make so user only can delete self
}
