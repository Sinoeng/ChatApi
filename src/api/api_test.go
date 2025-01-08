package api

import (
	"bytes"
	"os"
	"primary/database"
	"primary/pubsub"
	"testing"

	jwtMiddleware "primary/api/middleware/jwt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var testEngine *gin.Engine

func TestMain(m *testing.M) {
	initDB()

	sendCh := make(chan pubsub.Message, 100)
	go func() {
		for {
			_ = <-sendCh
		}
	}()

    testEngine = InitRouter(db, sendCh)
    go testEngine.Run()

	m.Run()

	db.ResetTestDBApi()
}

func initDB() {
	var err error
	db, err = database.InitDatabase(os.Getenv("MYSQL_DATABASE") + "Api")
	if err != nil {
		println("Failed to init db")
		panic(err)
	}
}

func resetDB() {
	db.ResetTestDBApi()
	initDB()
}

func checkErr(t *testing.T, desc string, err error) {
    if err != nil {
        t.Fatalf(desc + ". Err: %s\n", err.Error())
    }
}

func generateBody(body string) *bytes.Buffer {
    return bytes.NewBuffer([]byte(body))
}

func setClaims(username string, userID uint) *gin.Context {
	ctx := gin.Context{}
	claims := jwtMiddleware.MyCustomClaims{
		Username: username,
		Userid:   userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "ChatApi",
		},
	}
	ctx.Set("claims", &claims)
	return &ctx
}
