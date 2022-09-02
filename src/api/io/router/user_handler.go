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
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   err.Error(),
			"message": "error reading request body!",
		})
		return
	}

	newUser, err := io.DeserializeUser(jsonBody)
	if err != nil {
		log.Printf("error unmarshalling user body: %v", err)
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   err.Error(),
			"message": "error reading user body!",
		})
		return
	}

	userManager := users.NewUserManager(env.UserPersister)
	user, err := userManager.Register(newUser)
	if err != nil {
		log.Printf("%v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "error inserting new user!",
		})
		return
	}

	c.JSON(http.StatusOK, user)
}

func UpdateUser(c *gin.Context, env environment.Env) {
	jsonBody, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Printf("error reading request body: %v", err)
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   err.Error(),
			"message": "error reading request body!",
		})
		return
	}

	user, err := io.DeserializeUser(jsonBody)
	if err != nil {
		log.Printf("error unmarshalling user body: %v", err)
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   err.Error(),
			"message": "error reading user body!",
		})
		return
	}

	userManager := users.NewUserManager(env.UserPersister)
	updatedUser, err := userManager.Modify(user)
	if err != nil {
		log.Printf("error updating user with ID %s: %v ", user.ID, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "error updating user!",
		})
		return
	}

	c.JSON(http.StatusOK, updatedUser)
}
