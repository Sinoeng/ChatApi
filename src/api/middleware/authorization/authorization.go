package authorization

import (
	"net/http"
	"primary/database"
	"primary/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

func AuthorizeMiddleware(c *gin.Context, db *database.ChatApiDB, checks ...func(*gin.Context, *database.ChatApiDB) bool) {
	for _, check := range checks {
		if check(c, db) {
			c.Next()
			return
		}
	}
	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
	return
}

func CheckGlobalAdmin(c *gin.Context, db *database.ChatApiDB) bool {
	claims, err := utils.GetClaims(c)
	if err != nil {
		return false
	}

	user, err := db.GetUserByID(claims.Userid)
	if err != nil {
		return false
	}

	return user.Admin
}

func CheckServerAdmin(c *gin.Context, db *database.ChatApiDB) bool {
	claims, err := utils.GetClaims(c)
	if err != nil {
		return false
	}

	serverID, err := strconv.ParseUint(c.Param("serverid"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error:": "serverid is not a number"})
		return false
	}

	userServer, err := db.GetUserServerByIDs(claims.Userid, uint(serverID))
	if err != nil {
		return false
	}

	if userServer.Role == database.ROLE_SERVER_ADMIN {
		return true
	}

	return false
}

func CheckSameUser(c *gin.Context, db *database.ChatApiDB) bool {
	claims, err := utils.GetClaims(c)
	if err != nil {
		return false
	}

	userID, err := strconv.ParseUint(c.Param("userid"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error:": "userid is not a number"})
		return false
	}

	if claims.Userid == uint(userID) {
		return true
	}
	return false
}

func CheckServerMember(c *gin.Context, db *database.ChatApiDB) bool {
	claims, err := utils.GetClaims(c)
	if err != nil {
		return false
	}

	serverID, err := strconv.ParseUint(c.Param("serverid"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error:": "serverid is not a number"})
		return false
	}

	_, err = db.GetUserServerByIDs(claims.Userid, uint(serverID))
	if err != nil {
		return false
	}

	return true
}
