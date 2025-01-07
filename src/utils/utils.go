package utils

import (
	"errors"
	"log"
	"net/http"
	"os"
	jwtMiddleware "primary/api/middleware/jwt"
	"primary/database"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func GetClaims(c *gin.Context) (*jwtMiddleware.MyCustomClaims, error) {
	claimsAny, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "faulty token"})
		return nil, errors.New("CantGetClaimsError")
	}
	claims := claimsAny.(*jwtMiddleware.MyCustomClaims)
	return claims, nil
}

func CreateDefaultAdmin(db *database.ChatApiDB) {
    count, err := db.GetAdminUserCount()
    if err != nil {
        log.Printf("Failed to get how number of admin users. Err %s\n", err.Error())
        return
    }
    if count > 0 {
        return
    }
    username := os.Getenv("DEFAULT_USERNAME")
    password, err := bcrypt.GenerateFromPassword([]byte(os.Getenv("DEFAULT_PASSWORD")), 11)
    if err != nil {
        log.Printf("Failed to generate the default admin password. Err: %s\n", err.Error())
        return
    }
    _, err = db.InsertNewUserAsAdmin(username, string(password))
    if err != nil {
        log.Printf("Failed to create default admin user. Err: %s\n", err.Error())
    }
    log.Println("Default admin user created")
}
