package jwt

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type MyCustomClaims struct {
	Username string `json:"username"`
	Userid   uint   `json:"userid"`
	jwt.RegisteredClaims
}

func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		unparsedToken := c.Request.Header.Get("Authorization")
		fmt.Println("unparsed token:", unparsedToken)

		token, err := jwt.ParseWithClaims(unparsedToken, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_KEY")), nil
		})

		if err != nil {
			fmt.Println(err)
			return
		}

		if claims, ok := token.Claims.(*MyCustomClaims); ok {
			fmt.Println("JWT VALIDATED")
			fmt.Print(claims.Username)
			c.Set("claims", claims)
			c.Next()
		} else {
			fmt.Println(err)
			return
		}
	}
}
