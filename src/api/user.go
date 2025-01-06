package api

import (
	"net/http"
	"primary/api/middleware/authorization"
	"strconv"

	"github.com/gin-gonic/gin"
)

func getServersHandler(c *gin.Context) {
	userid64, err := strconv.ParseUint(c.Param("userid"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userid is not a number"})
		return
	}

	servers, err := db.GetAllServersByUserID(uint(userid64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error:": "could not find any servers"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": servers})
}

func removeUserHandler(c *gin.Context) { // TODO: add authentication
	userid64, err := strconv.ParseUint(c.Param("userid"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userid is not a number"})
		return
	}

	if err := db.DeleteUserByID(uint(userid64)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Could not delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "deleted user " + strconv.FormatInt(int64(userid64), 10)})
}

func AddUserRoutes(grp *gin.RouterGroup) {
	grp.DELETE("/:userid", func(c *gin.Context) {
		authorization.AuthorizeMiddleware(c, db, authorization.CheckGlobalAdmin, authorization.CheckSameUser)
	}, removeUserHandler) //TODO make so user only can delete self
	grp.GET("/:userid/servers", func(c *gin.Context) {
		authorization.AuthorizeMiddleware(c, db, authorization.CheckGlobalAdmin, authorization.CheckSameUser)
	}, getServersHandler)
}
