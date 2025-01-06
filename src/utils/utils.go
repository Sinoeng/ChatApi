package utils

import (
	"errors"
	"net/http"
	jwtMiddleware "primary/api/middleware/jwt"

	"github.com/gin-gonic/gin"
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
