package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// PingExample godoc
// @Summary ping example
// @Schemes
// @Description do ping
// @Tags example
// @Accept json
// @Produce json
// @Success 200 {string} Helloworld
// @Router /ping/ [get]
func PingHandler(c *gin.Context) {
	c.JSON(http.StatusOK, "pong")
}
