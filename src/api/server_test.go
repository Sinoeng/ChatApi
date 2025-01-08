package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
)


func TestNewServerHandler(t *testing.T) {
    resetDB()

    body := generateBody(`{
        "name": "servername"
        }`)
    req, err := http.NewRequest("POST", "/v1/protected/server/new", body)
    checkErr(t, "Failed to create http request", err)
    req.Header.Add("Content-Type", "application/json")

    w := httptest.NewRecorder()
    testEngine.ServeHTTP(w, req)

    if w.Code != http.StatusInternalServerError {
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

