package router

import "github.com/gin-gonic/gin"

const (
	pingPath = "/ping"
)

var (
	Router *gin.Engine
)
