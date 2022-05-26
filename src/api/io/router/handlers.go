package router

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type HandlerService struct {
}

func RegisterNewUser(c *gin.Context) {
	c.String(http.StatusOK, "user register correctly!")
}
