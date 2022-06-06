package router

import (
	"be-tpms/src/api/configuration"
	"be-tpms/src/api/environment"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
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

func Run(env environment.Env, port string) error {
	SetupRunEnv(env)
	err := Router.Run(":" + port)
	return err
}

func SetupRunEnv(env environment.Env) {
	health := HealthChecker{}
	log.Print("[package:router] Configuring routes...")
	ConfigureRoute(env, health)
	log.Printf(fmt.Sprintf("[package:router] Listening on routes: %s", pingPath))
}

func ConfigureRoute(env environment.Env, health HealthChecker) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = ioutil.Discard
	Router = gin.Default()
	mapRoutes(Router, health, env)
	return Router
}

func mapRoutes(r *gin.Engine, health HealthChecker, env environment.Env) {
	r.GET(pingPath, health.PingHandler)
	r.POST(registerUserPath, RegisterNewUser)
}

type HealthChecker struct{}

func (h HealthChecker) PingHandler(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}

func CreateRestClientConfig(profile string) *resty.Client {
	restClient := resty.New()
	if profile == configuration.Test {
		restClient.Header.Add("test", "true")
	}
	return restClient
}
