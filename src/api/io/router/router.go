package router

import (
	_ "be-tpms/docs"
	"be-tpms/middleware"
	"be-tpms/src/api/environment"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	basePath    = "/api/v1"
	pingPath    = "/ping"
	userPath    = "/user"
	dogPath     = "/dog"
	imgPath     = "/img"
	swaggerPath = "/swagger"
)

var (
	router *gin.Engine
)

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

func CreateRestClientConfig() *http.Client {
	return http.DefaultClient
}

func mapDogRoutes(env environment.Env) {
	dogRouter := router.Group(basePath + dogPath)
	dogRouter.POST("", func(context *gin.Context) {
		if !validUser(context) {
			return
		}
		RegisterNewDog(context, env)
	})
	dogRouter.GET("/:id", func(context *gin.Context) {
		if !validUser(context) {
			return
		}
		GetDog(context, env)
	})
	dogRouter.DELETE("/:id", func(context *gin.Context) {
		if !validUser(context) {
			return
		}
		DeleteDog(context, env)
	})
	dogRouter.PATCH("", func(context *gin.Context) {
		if !validUser(context) {
			return
		}
		UpdateDog(context, env)
	})
	dogRouter.PUT("/found", func(context *gin.Context) {
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
	dogRouter.POST("/possible", func(context *gin.Context) {
		if !validUser(context) {
			return
		}
		PossibleMatch(context, env)
	})
	dogRouter.PUT("/possible", func(context *gin.Context) {
		if !validUser(context) {
			return
		}
		AckPossibleDog(context, env)
	})
	dogRouter.DELETE("/possible", func(context *gin.Context) {
		if !validUser(context) {
			return
		}
		RejectPossibleDog(context, env)
	})
	dogRouter.GET("/:id/possible", func(context *gin.Context) {
		if !validUser(context) {
			return
		}
		GetPossibleMatchingDog(context, env)
	})
	dogRouter.GET("/:id/prediction", func(context *gin.Context) {
		if !validUser(context) {
			return
		}
		GetSimilarDogPrediction(context, env)
	})
}

func mapUserRoutes(env environment.Env) {
	userRouter := router.Group(basePath + userPath)
	userRouter.POST("", func(context *gin.Context) {
		if !validUser(context) {
			return
		}
		RegisterNewUser(context, env)
	})
	userRouter.PATCH("", func(context *gin.Context) {
		if !validUser(context) {
			return
		}
		UpdateUser(context, env)
	})
	userRouter.GET("/dog", func(context *gin.Context) {
		if !validUser(context) {
			return
		}
		GetUserDogs(context, env)
	})
	userRouter.PUT("/fcmtoken", func(context *gin.Context) {
		if !validUser(context) {
			return
		}
		UpdateFCMToken(context, env)
	})
	userRouter.GET("", func(context *gin.Context) {
		if !validUser(context) {
			return
		}
		GetUser(context, env)
	})
	userRouter.GET("/:id", func(context *gin.Context) {
		if !validUser(context) {
			return
		}
		GetUserContactInfo(context, env)
	})
	userRouter.POST("/:id/notification", func(context *gin.Context) {
		SendNotif(context, env)
	})
}

func mapImgsRoutes(env environment.Env) {
	imgsRouter := router.Group(basePath + imgPath)
	imgsRouter.POST("", func(context *gin.Context) {
		if !validUser(context) {
			return
		}
		AddImg(context, env)
	})
}

func mapPingRoutes() {
	router.GET(basePath+pingPath, Ping)
}

func mapSwaggerRoutes() {
	swaggerRouter := router.Group(basePath + swaggerPath)
	swaggerRouter.GET("/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}

func validUser(c *gin.Context) bool {
	middleware.AuthMiddleware(c)
	return c.GetHeader("x-user-id") != ""
}
