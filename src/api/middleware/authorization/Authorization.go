package authorization

import (
	"net/http"
	"primary/database"

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
	/*claims, err := utils.GetClaims(c)
	if err != nil {
		return false
	}

	user, err := db.GetUserByID(claims.Userid)
	if err != nil {
		return false
	}*/

	return false
}

func CheckServerAdmin(c *gin.Context, db *database.ChatApiDB) bool {
	return true
}

func CheckSameUser(c *gin.Context, db *database.ChatApiDB) bool {
	return true
}

func CheckServerMember(c *gin.Context, db *database.ChatApiDB) bool {
	return true
}
