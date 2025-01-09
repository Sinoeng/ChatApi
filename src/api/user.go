package api

import (
	"net/http"
	"primary/api/middleware/authorization"
	"primary/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type NewPassword struct {
	Password string `form:"password" json:"password" xml:"password" binding:"required"`
}

// get_servers godoc
// @Summary get servers of a user
// @Schemes
// @Description get the servers user with userid is a member of
// @Tags auth user
// @Param Authorization header string true "jwt token"
// @Produce json
// @Success 200
// @Router /protected/user/:userid/servers [get]
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

// remove_user godoc
// @Summary delete a user if you are that user or an admin
// @Schemes
// @Description get the servers user with userid is a member of
// @Param Authorization header string true "jwt token"
// @Tags auth user
// @Produce json
// @Success 200
// @Router /protected/user/:userid [delete]
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

// change_password godoc
// @Summary change user password
// @Schemes
// @Param Authorization header string true "jwt token"
// @Param password body string true "new password"
// @Tags auth user
// @Accept json
// @Produce json
// @Success 200
// @Router /protected/user/:userid/changepassword [patch]
func changePasswordHandler(c *gin.Context) {
	claims, err := utils.GetClaims(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "faulty token"})
		return
	}

	userID := claims.Userid

	var newPassword NewPassword
	if err := c.Bind(&newPassword); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "could not read password"})
		return
	}
	bytes, err := bcrypt.GenerateFromPassword([]byte(newPassword.Password), 11)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "could not hash password"})
		return
	}

	if err = db.ChangeUserPasswordByID(userID, string(bytes)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not update password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "password updated"})
}

func AddUserRoutes(grp *gin.RouterGroup) {
	grp.DELETE("/:userid", func(c *gin.Context) {
		authorization.AuthorizeMiddleware(c, db, authorization.CheckGlobalAdmin, authorization.CheckSameUser)
	}, removeUserHandler) //TODO make so user only can delete self
	grp.GET("/:userid/servers", func(c *gin.Context) {
		authorization.AuthorizeMiddleware(c, db, authorization.CheckGlobalAdmin, authorization.CheckSameUser)
	}, getServersHandler)
	grp.PATCH("/:userid/changepassword", func(c *gin.Context) {
		authorization.AuthorizeMiddleware(c, db, authorization.CheckGlobalAdmin, authorization.CheckSameUser)
	}, changePasswordHandler)
}
