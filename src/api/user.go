package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type User struct {
	Name     string `form:"name" json:"name" xml:"name"  binding:"required"`
	Password string `form:"password" json:"password" xml:"password" binding:"required"`
}

// TODO update status codes and look over error handling for ALL functions

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
	tokenName := claims["user"].(string)

	if requestName != tokenName {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "cannot delete users other than yourself"})
		return
	}

	dbUsr, err := db.GetUserByUsername(tokenName)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
}
