package router

import (
	"be-tpms/src/api/configuration"
	"be-tpms/src/api/environment"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"io/ioutil"
	"log"
)

const (
	pingPath = "/ping"
	userPath = "/user"
	dogPath  = "/dog"
)

var (
	router *gin.Engine
)

func Run(env environment.Env, port string) error {
	SetupRunEnv(env)
	err := router.Run(":" + port)
	return err
}

func mapHandlers(env environment.Env) {
	mapPingRoutes()
	mapUserRoutes(env)
}

func SetupRunEnv(env environment.Env) {
	log.Print("[package:router] Configuring routes...")
	ConfigureRoute(env)
	log.Printf(fmt.Sprintf("[package:router] Listening on routes: %s, %s", pingPath, userPath))
}

func ConfigureRoute(env environment.Env) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = ioutil.Discard
	router = gin.Default()
	mapHandlers(env)
	return router
}

func CreateRestClientConfig(profile string) *resty.Client {
	restClient := resty.New()
	if profile == configuration.Test {
		restClient.Header.Add("test", "true")
	}
	return restClient
}

func mapUserRoutes(env environment.Env) {
	router.POST(userPath, func(context *gin.Context) {
		RegisterNewUser(context, env)
	})
	router.PATCH(userPath, func(context *gin.Context) {
		UpdateUser(context, env)
	})
	router.POST(dogPath, func(context *gin.Context) {
		RegisterNewDog(context, env)
	})
	router.PATCH(dogPath, func(context *gin.Context) {
		UpdateDog(context, env)
	})
}

func mapPingRoutes() {
	router.GET(pingPath, PingHandler)
}
