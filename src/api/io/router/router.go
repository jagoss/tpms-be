package router

import (
	"be-tpms/src/api/environment"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	pingPath         = "/ping"
	registerUserPath = "/user"
)

var (
	Router *gin.Engine
)

func Run(env environment.Environment, port string) error {
	SetupRunEnv(env)
	err := Router.Run(":" + port)
	return err
}

func SetupRunEnv(env environment.Environment) {
	health := HealthChecker{}
	log.Print("[package:router] Configuring routes...")
	ConfigureRoute(env, health)
	log.Printf(fmt.Sprintf("[package:router] Listening on routes: %s", pingPath))
}

func ConfigureRoute(env environment.Environment, health HealthChecker) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = ioutil.Discard
	Router = gin.Default()
	mapRoutes(Router, health, env)
	return Router
}

func mapRoutes(r *gin.Engine, health HealthChecker, env environment.Environment) {
	r.GET(pingPath, health.PingHandler)
	r.POST(registerUserPath, RegisterNewUser)
}

type HealthChecker struct{}

func (h HealthChecker) PingHandler(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}
