package router

import (
	"be-tpms/middleware"
	"be-tpms/src/api/configuration"
	"be-tpms/src/api/environment"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"io/ioutil"
	"log"
)

const (
	pingPath = "/ping"
	userPath = "/user"
	dogPath  = "/dog"
	imgPath  = "/img"
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
	mapDogRoutes(env)
	mapImgsRoutes(env)
}

func SetupRunEnv(env environment.Env) {
	log.Print("[package:router] Configuring routes...")
	configureRoute(env)
	log.Printf("[package:router] Listening on routes: %s, %s, %s, %s", pingPath, userPath, dogPath, imgPath)
}

func configureRoute(env environment.Env) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = ioutil.Discard
	router = gin.Default()
	router.Use(func(c *gin.Context) {
		c.Set("firebaseAuth", &env.FirebaseAuth)
	})
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

func mapDogRoutes(env environment.Env) {
	router.POST(dogPath, func(context *gin.Context) {
		if !validUser(context) {
			return
		}
		RegisterNewDog(context, env)
	})
	router.PATCH(dogPath, func(context *gin.Context) {
		if !validUser(context) {
			return
		}
		UpdateDog(context, env)
	})
	router.PATCH(dogPath+"/found", func(context *gin.Context) {
		if !validUser(context) {
			return
		}
		DogReUnited(context, env)
	})
	router.GET(dogPath+"/missing", func(context *gin.Context) {
		if !validUser(context) {
			return
		}
		GetMissingDogsList(context, env)
	})
}

func mapUserRoutes(env environment.Env) {
	router.POST(userPath, func(context *gin.Context) {
		middleware.AuthMiddleware(context)
		RegisterNewUser(context, env)
	})
	router.PATCH(userPath, func(context *gin.Context) {
		if !validUser(context) {
			return
		}
		UpdateUser(context, env)
	})
}

func mapImgsRoutes(env environment.Env) {
	router.POST(imgPath, func(context *gin.Context) {
		//if !validUser(context) {
		//	return
		//}
		AddImg(context, env)
	})
}

func mapPingRoutes() {
	router.GET(pingPath, PingHandler)
}

func validUser(c *gin.Context) bool {
	middleware.AuthMiddleware(c)
	_, exists := c.Get("UUID")
	return exists
}
