package router

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// swagger:route GET /ping ping
// Ping server to check if it's up and running
//
// security:
// - apiKey: []
// responses:
//  200: pong
func PingHandler(c *gin.Context) {
	c.JSON(http.StatusOK, "pong")
}
