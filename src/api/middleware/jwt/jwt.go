package jwt

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {

		unparsedToken := c.Request.Header.Get("token")
		fmt.Println("unparsed token:", unparsedToken)

		token, err := jwt.Parse(unparsedToken, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_KEY")), nil
		})

		if err != nil {
			fmt.Println(err)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			fmt.Print(claims["username"])
		} else {
			fmt.Println(err)
			return
		}

		fmt.Println("JWT VALIDATED")
		c.Set("token", token)
		c.Next()
	}
}
