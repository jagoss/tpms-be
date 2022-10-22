package router

import (
	"be-tpms/src/api/environment"
	"be-tpms/src/api/io"
	"be-tpms/src/api/usecases/dogs"
	"be-tpms/src/api/usecases/users"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
)

type UserHandler struct {
}

// RegisterNewUser godoc
// @Summary Register new user
// @Schemes
// @Description Register new user
// @Tags        user
// @Accept      json
// @Produce     json
// @Param		user body model.User false "new user"
// @Success     200 {object} model.User
// @Failure		422 {object} object{error=string, message=string}
// @Failure		500 {object} object{error=string, message=string}
// @Router      /user [post]
func RegisterNewUser(c *gin.Context, env environment.Env) {
	jsonBody, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Printf("error reading request body: %v", err)
		c.JSON(http.StatusUnprocessableEntity, map[string]string{
			"error":   err.Error(),
			"message": "error reading request body!",
		})
		return
	}

	newUser, err := io.DeserializeUser(jsonBody)
	newUser.ID = c.Request.Header.Get("user_id")
	if err != nil {
		log.Printf("error unmarshalling user body: %v", err)
		c.JSON(http.StatusUnprocessableEntity, map[string]string{
			"error":   err.Error(),
			"message": "error reading user body!",
		})
		return
	}

	userManager := users.NewUserManager(env.UserPersister)
	user, err := userManager.Register(newUser)
	if err != nil {
		log.Printf("%v", err)
		c.JSON(http.StatusInternalServerError, map[string]string{
			"error":   err.Error(),
			"message": "error inserting new user!",
		})
		return
	}

	c.JSON(http.StatusOK, user)
}

// UpdateUser godoc
// @Summary Updates user
// @Schemes
// @Description Updates existing user
// @Tags        user
// @Accept      json
// @Produce     json
// @Param		user body model.User true "user to update"
// @Success     200 {object} model.User
// @Failure		422 {object} object{error=string, message=string}
// @Failure		500 {object} object{error=string, message=string}
// @Router      /user [patch]
func UpdateUser(c *gin.Context, env environment.Env) {
	jsonBody, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Printf("error reading request body: %v", err)
		c.JSON(http.StatusUnprocessableEntity, map[string]string{
			"error":   err.Error(),
			"message": "error reading request body!",
		})
		return
	}

	user, err := io.DeserializeUser(jsonBody)
	user.ID = c.Request.Header.Get("user_id")
	if err != nil {
		log.Printf("error unmarshalling user body: %v", err)
		c.JSON(http.StatusUnprocessableEntity, map[string]string{
			"error":   err.Error(),
			"message": "error reading user body!",
		})
		return
	}

	userManager := users.NewUserManager(env.UserPersister)
	updatedUser, err := userManager.Modify(user)
	if err != nil {
		log.Printf("error updating user with ID %s: %v ", user.ID, err)
		c.JSON(http.StatusInternalServerError, map[string]string{
			"error":   err.Error(),
			"message": "error updating user!",
		})
		return
	}

	c.JSON(http.StatusOK, updatedUser)
}

// GetUserDogs godoc
// @Summary Get all user dogs
// @Schemes
// @Description Gets 2 lists of dogs, one with dogs owned by de user and another with dogs found by the user
// @Tags        user
// @Accept      json
// @Produce     json
// @Param		user body model.User true "user to update"
// @Success     200 {object} object{ownedDogs=[]model.DogResponse, foundDogs=[]model.DogResponse}
// @Failure		500 {object} object{error=string, message=string}
// @Router      /user/dog [get]
func GetUserDogs(c *gin.Context, env environment.Env) {
	userID := c.GetHeader("user_id")
	dogManager := dogs.NewDogManager(env.DogPersister, env.Storage)
	userOwnedDogs, foundDogs, err := dogManager.GetAllUserDogs(userID)
	if err != nil {
		msg := fmt.Sprintf("Error getting dogs for user %s", userID)
		log.Printf(msg)
		c.JSON(http.StatusInternalServerError, map[string]string{
			"error":   err.Error(),
			"message": msg,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ownedDogs": io.MapToDogResponseList(userOwnedDogs, env.Storage),
		"foundDogs": io.MapToDogResponseList(foundDogs, env.Storage),
	})
}

// UpdateFCMToken godoc
// @Summary Updates FCM token for given user
// @Schemes
// @Description Updates FCM token to allow server to send direct notifications to users at any time
// @Tags        user
// @Accept      json
// @Produce     json
// @Param		user body object{token=string} true "FCM token"
// @Success     200 {object} object{result=string}
// @Failure		422 {object} object{error=string, message=string}
// @Failure		500 {object} object{error=string, message=string}
// @Router      /user/fcmtoken [put]
func UpdateFCMToken(c *gin.Context, env environment.Env) {
	userID := c.GetHeader("user_id")
	jsonBody, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Printf("error reading request body: %v", err)
		c.JSON(http.StatusUnprocessableEntity, map[string]string{
			"error":   err.Error(),
			"message": "error reading request body!",
		})
		return
	}

	var body map[string]string
	err = json.Unmarshal(jsonBody, &body)
	if err != nil {
		log.Printf("error reading request body: %v", err)
		c.JSON(http.StatusInternalServerError, map[string]string{
			"error":   err.Error(),
			"message": "error unmarshalling request body!",
		})
		return
	}

	userManager := users.NewUserManager(env.UserPersister)
	err = userManager.UpdateFCMToken(userID, body["token"])
	if err != nil {
		msg := fmt.Sprintf("error saving fcm token for user %s", userID)
		log.Printf(msg+": %v", err)
		c.JSON(http.StatusInternalServerError, map[string]string{
			"error":   err.Error(),
			"message": msg,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": "OK",
	})
}
