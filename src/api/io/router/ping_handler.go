package router

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// Ping godoc
// @Summary Ping
// @Schemes
// @Description do ping
// @Tags        ping
// @Accept      json
// @Produce     json
// @Success     200 {string} pong
// @Router      /ping [get]
func Ping(c *gin.Context) {
	c.JSON(http.StatusOK, "pong")
}
