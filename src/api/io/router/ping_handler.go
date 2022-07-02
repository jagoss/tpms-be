package router

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func PingHandler(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}
