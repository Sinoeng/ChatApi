package api

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func LoginHandler(c *gin.Context) { // issue jwt
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

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user": creds.Name,
		"iat":  time.Now(),
		"exp":  time.Now().Add(7 * 24 * time.Hour), //TODO: add user role
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_KEY")))

	c.JSON(http.StatusOK, gin.H{"status": "Logged in successfully", "token": tokenString}) //TODO: issue jwt
}

func NewHandler(c *gin.Context) { //TODO: add email
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

	if usr.Email == "" {
		if _, err := db.InsertNewUser(usr.Name, string(bytes)); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	} else {
		if _, err := db.InsertNewUserWEmail(usr.Name, string(bytes), usr.Email); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"status": "user created successfully"})
}
