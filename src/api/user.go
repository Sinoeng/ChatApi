package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Name     string `form:"name" json:"name" xml:"name"  binding:"required"`
	Password string `form:"password" json:"password" xml:"password" binding:"required"`
}

// TODO update status codes and look over error handling for ALL functions
func new(c *gin.Context) { //TODO: add email
	var usr User

	if err := c.Bind(&usr); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	bytes, err := bcrypt.GenerateFromPassword([]byte(usr.Password), 11)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "could not hash password"})
		return
	}

	if err := db.InsertNewUser(usr.Name, string(bytes)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "user created successfully"})
}

func verify(c *gin.Context) { // verify jwt

}

func login(c *gin.Context) { // issue jwt
	var creds User

	if err := c.Bind(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := db.GetUserByUsername(creds.Name)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Logged in successfully"}) //TODO: issue jwt
}

func remove(c *gin.Context) { // TODO: add authentication
	var reqUsr User

	if err := c.Bind(&reqUsr); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	dbUsr, err := db.GetUserByUsername(reqUsr.Name)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.DeleteUserByID(dbUsr.ID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

}

func AddUserRoutes(grp *gin.RouterGroup) {

	user := grp.Group("/user")

	user.POST("/new", new)
	user.POST("/login", login)
	user.DELETE("/:user", remove)

}
