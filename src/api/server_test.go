package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"primary/utils"
	"testing"

	"github.com/gin-gonic/gin"
)


func TestNewServerHandler(t *testing.T) {
    resetDB()
    utils.CreateDefaultAdmin(db)
    body := generateBody(fmt.Sprintf(`{
        "name": "%s",
        "password": "%s"
        }`, os.Getenv("DEFAULT_USERNAME"), os.Getenv("DEFAULT_PASSWORD")))
    w := httptest.NewRecorder()
    req, err := http.NewRequest("POST", "/v1/user/login", body)
    checkErr(t, "Failed to create http request", err)
    req.Header.Add("Content-Type", "application/json")

    testEngine.ServeHTTP(w, req)
    if w.Code != http.StatusOK {
        t.Fatalf("Faulty error code %d\nBody: %s\n", w.Code, w.Body.String())
    }
    var js, data gin.H
    err = json.Unmarshal(w.Body.Bytes(), &js)
    checkErr(t, "Failed to parse the login response", err)
    data, ok := js["data"].(map[string]interface{})
    if !ok {
        t.Fatal("Failed to extract 'data' field")
    }
    token := data["token"].(string)

    body = generateBody(`{
        "name": "servername"
        }`)
    req, err = http.NewRequest("POST", "/v1/protected/server/new", body)
    checkErr(t, "Failed to create http request 2", err)
    req.Header.Add("Content-Type", "application/json")
    req.Header.Add("Authorization", token)

    testEngine.ServeHTTP(w, req)

    if w.Code != http.StatusOK {
        t.Fatal("Faulty status code")
    }

    //
    // username := "jack"
    // var userID uint = 2
    // ctx := setClaims(username, userID)
    //
    // ctx.Request = req
    // log.Printf("resp %v\n", ctx.Request.Response)
    // ctx.Request.Response = &http.Response{}
    // w := httptest.NewRecorder()
    // ctx, _ = gin.CreateTestContext(w)
    // ctx.JSON(500, "sd")
    // newServerHandler(ctx)
    //
    // if ctx.Request.Response.Status != "200" {
    //     t.Fatal("Got a faulty status code")
    // }
    // t.Fatalf("StatusCode %d\n", ctx.Request.Response.StatusCode)
}

