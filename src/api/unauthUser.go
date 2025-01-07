package api

import (
	"fmt"
	"net/http"
	"os"
	jwtMiddleware "primary/api/middleware/jwt"
	"primary/pubsub"
	"regexp"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Name     string `form:"name" json:"name" xml:"name"  binding:"required"`
	Password string `form:"password" json:"password" xml:"password" binding:"required"`
	Email    string `form:"email" json:"email" xml:"email"`
}

var emailRegex = regexp.MustCompile(`^[\w-\.]+@([\w-]+\.)+[\w-]{2,4}$`)

func loginHandler(c *gin.Context) { // issue jwt
	var creds User

	if err := c.Bind(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "could not extract user from request"})
		return
	}

	user, err := db.GetUserPasswordByUsername(creds.Name)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "password or username icorrect"})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "password or username icorrect"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtMiddleware.MyCustomClaims{ //TODO: make issue jwt function
		Username: creds.Name,
		Userid:   user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "ChatApi",
		},
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_KEY")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not issue jwt"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Logged in successfully", "token": tokenString, "userid": user.ID}) //TODO: issue jwt
}

func newHandler(c *gin.Context) {
	var usr User

	if err := c.Bind(&usr); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "could not get user information"})
		return
	}

	bytes, err := bcrypt.GenerateFromPassword([]byte(usr.Password), 11)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "could not hash password"})
		return
	}

	if usr.Email == "" {
		if _, err := db.InsertNewUser(usr.Name, string(bytes)); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "could not create user"})
			return
		}
	} else {
		if isValidEmail := emailRegex.MatchString(usr.Email); !isValidEmail {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email address"})
			return
		}

		if _, err := db.InsertNewUserWEmail(usr.Name, string(bytes), usr.Email); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "could not create user"})
			return
		}
		toPSCh <- pubsub.Message{Payload: fmt.Sprintf("A new user %s has been made.", usr.Name), Channel: os.Getenv("EMAIL_CHANNEL"), Receiver: usr.Email, Subject: "Welcome to ChatApi"}
	}

	c.JSON(http.StatusOK, gin.H{"status": "user created successfully"})
}

func AddUnauthUserRoutes(grp *gin.RouterGroup) {
	grp.POST("/login", loginHandler)
	grp.POST("/newuser", newHandler)
}
