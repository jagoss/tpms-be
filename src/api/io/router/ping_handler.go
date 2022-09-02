package router

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// PingExample godoc
// @Summary ping example
// @Schemes
// @Description do ping
// @Tags        example
// @Accept      json
// @Produce     json
// @Success     200 {string} pong
// @Router      /ping [get]
func PingHandler(c *gin.Context) {
	c.JSON(http.StatusOK, "pong")
}
