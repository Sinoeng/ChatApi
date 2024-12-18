package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type User struct {
	Name     string `form:"name" json:"name" xml:"name"  binding:"required"`
	Password string `form:"password" json:"password" xml:"password" binding:"required"`
	Email    string `form:"email" json:"email" xml:"email"`
}

// TODO update status codes and look over error handling for ALL functions

func GetServersHandler(c *gin.Context) {
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

	servers, err := db.GetAllServersByUserID(userid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error:": "could not find any servers"})
	}

	c.JSON(http.StatusOK, gin.H{"data": servers})
}

func removeHandler(c *gin.Context) { // TODO: add authentication
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

	requestName := c.Param("user")
	dbUsr, err := db.GetUserByUsername(requestName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	requestID := dbUsr.ID
	tokenID := uint(claims["userid"].(float64))

	if requestID != tokenID {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "cannot delete users other than yourself"})
		return
	}

	if err := db.DeleteUserByID(dbUsr.ID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "deleted user " + requestName})
}

func AddUserRoutes(grp *gin.RouterGroup) {
	user := grp.Group("/user")

	user.DELETE("/:user", removeHandler) //TODO make so user only can delete self
	user.GET("/servers", GetServersHandler)
}
