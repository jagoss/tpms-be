package router

import (
	"be-tpms/src/api/environment"
	"be-tpms/src/api/io"
	"be-tpms/src/api/usecases/users"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
)

type UserHandler struct {
}

func RegisterNewUser(c *gin.Context, env environment.Env) {
	jsonBody, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Printf("error reading request body: %v", err)
		c.String(http.StatusUnprocessableEntity, "error reading request body!")
	}
	newUser, err := io.DeserializeUser(jsonBody)
	if err != nil {
		log.Printf("error unmarshalling user body: %v", err)
		c.String(http.StatusUnprocessableEntity, "error reading user body!")
	}
	userManager := users.NewUserManager(&env.UserPersister)
	user, err := userManager.Register(newUser)
	if err != nil {
		log.Printf("%v", err)
		c.String(http.StatusInternalServerError, "error inserting new user")
	}

	c.JSON(http.StatusOK, user)
}

func UpdateUser(c *gin.Context, env environment.Env) {
	jsonBody, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Printf("error reading request body: %v", err)
		c.String(http.StatusUnprocessableEntity, "error reading request body!")
	}
	user, err := io.DeserializeUser(jsonBody)
	if err != nil {
		log.Printf("error unmarshalling user body: %v", err)
		c.String(http.StatusUnprocessableEntity, "error reading user body!")
	}
	userManager := users.NewUserManager(&env.UserPersister)
	updatedUser, err := userManager.Modify(user)
	if err != nil {
		log.Printf("error updating user %s: %v ", user.ID, err)
		c.String(http.StatusInternalServerError, "error updating user")
	}
	c.JSON(http.StatusOK, updatedUser)
}
