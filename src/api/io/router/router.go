package router

import (
	"be-tpms/middleware"
	"be-tpms/src/api/configuration"
	"be-tpms/src/api/environment"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"io/ioutil"
	"log"
)

const (
	pingPath    = "/ping"
	userPath    = "/user"
	dogPath     = "/dog"
	imgPath     = "/img"
	swaggerPath = "/swagger"
)

var (
	router *gin.Engine
)

//  Tu perro me suena back-end Api:
//   version: 0.0.1
//   title: TPMS-BE Api
//  Schemes: http, https
//  Host: localhost:8080
//  BasePath: /
//  Produces:
//    - application/json
//
// swagger:meta

func Run(env environment.Env, port string) error {
	SetupRunEnv(env)
	err := router.Run(":" + port)
	return err
}

//armar grupos
func mapHandlers(env environment.Env) {
	mapPingRoutes()
	mapSwaggerRoutes()
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
	dogRouter := router.Group(dogPath)
	dogRouter.POST("", func(context *gin.Context) {
		if !validUser(context) {
			return
		}
		RegisterNewDog(context, env)
	})
	dogRouter.PATCH("", func(context *gin.Context) {
		if !validUser(context) {
			return
		}
		UpdateDog(context, env)
	})
	dogRouter.PATCH("/found", func(context *gin.Context) {
		if !validUser(context) {
			return
		}
		DogReUnited(context, env)
	})
	dogRouter.GET("/missing", func(context *gin.Context) {
		if !validUser(context) {
			return
		}
		GetMissingDogsList(context, env)
	})
}

func mapUserRoutes(env environment.Env) {
	userRouter := router.Group(userPath)
	userRouter.POST("", func(context *gin.Context) {
		middleware.AuthMiddleware(context)
		RegisterNewUser(context, env)
	})
	userRouter.PATCH("", func(context *gin.Context) {
		if !validUser(context) {
			return
		}
		UpdateUser(context, env)
	})
}

func mapImgsRoutes(env environment.Env) {
	imgsRouter := router.Group(imgPath)
	imgsRouter.POST("", func(context *gin.Context) {
		//if !validUser(context) {
		//	return
		//}
		AddImg(context, env)
	})
}

func mapPingRoutes() {
	router.GET(pingPath, PingHandler)
}

func mapSwaggerRoutes() {
	swaggerRouter := router.Group(swaggerPath)
	swaggerRouter.GET("/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

func validUser(c *gin.Context) bool {
	middleware.AuthMiddleware(c)
	_, exists := c.Get("UUID")
	return exists
}
