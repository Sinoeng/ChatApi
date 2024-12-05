package routes

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
func ping(c *gin.Context) {
	c.JSON(http.StatusOK, "pong")
}

func AddPingRoutes(grp *gin.RouterGroup) {
	pingGroup := grp.Group("/ping")

	pingGroup.GET("/", ping)
}
