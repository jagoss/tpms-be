package router

import (
	"be-tpms/src/api/environment"
	"be-tpms/src/api/io"
	"be-tpms/src/api/usecases/dogs"
	"be-tpms/src/api/usecases/lostandfound"
	"be-tpms/src/api/usecases/users"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type DogHandler struct {
}

// RegisterNewDog godoc
// @Summary Register new dog
// @Schemes
// @Description Register new dog
// @Tags        dog
// @Accept      json
// @Produce     json
// @Param		dog body model.DogRequest false  "dog"
// @Success     200 {object} model.DogResponse
// @Failure		422 {object} map[string]any{error=string, message=string}
// @Failure		500 {object} map[string]any{error=string, message=string}
// @Router      /dog [post]
func RegisterNewDog(c *gin.Context, env environment.Env) {
	jsonBody, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Printf("error reading request body: %v", err)
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   err.Error(),
			"message": "error reading request body!",
		})
		return
	}
	reqDog, err := io.DeserializeDog(jsonBody)
	if err != nil {
		log.Printf("error unmarshalling dog body: %v", err)
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   err.Error(),
			"message": "error reading dog body!",
		})
		return
	}
	dog, imgs := io.MapFromDogRequest(reqDog)
	if dog == nil && imgs == nil {
		log.Printf("error unmarshalling dog body: %v", err)
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   err.Error(),
			"message": "error parsing dogRequest to dog!",
		})
		return
	}

	userManager := users.NewUserManager(env.UserPersister)
	dogManager := dogs.NewDogManager(env.DogPersister, env.Storage)
	dog, err = dogManager.Register(dog, imgs, userManager)
	if err != nil {
		log.Printf("%v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": fmt.Sprintf("error inserting new dog: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, io.MapToDogResponse(dog, env.Storage))
}

// UpdateDog godoc
// @Summary Updates dog
// @Schemes
// @Description Updates existing dog
// @Tags        dog
// @Accept      json
// @Produce     json
// @Param		dog body model.DogRequest false  "dog"
// @Success     200 {object} model.DogResponse
// @Failure		422 {object} map[string]string{error=string, message=string}
// @Failure		500 {object} map[string]string{error=string, message=string}
// @Router      /dog [patch]
func UpdateDog(c *gin.Context, env environment.Env) {
	jsonBody, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Printf("error reading request body: %v", err)
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   err.Error(),
			"message": "error reading request body!",
		})
		return
	}

	reqDog, err := io.DeserializeDog(jsonBody)
	if err != nil {
		log.Printf("error unmarshalling dog body: %v", err)
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   err.Error(),
			"message": "error reading dpg body!",
		})
		return
	}

	dog, imgs := io.MapFromDogRequest(reqDog)
	if dog == nil && imgs == nil {
		log.Printf("error unmarshalling dog body: %v", err)
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   err.Error(),
			"message": "error parsing dogRequest to dog!",
		})
		return
	}

	userManager := users.NewUserManager(env.UserPersister)
	dogManager := dogs.NewDogManager(env.DogPersister, env.Storage)
	updatedDog, err := dogManager.Modify(dog, imgs, userManager)
	if err != nil {
		log.Printf("error updating dog with ID %s: %v ", reqDog.ID, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "error updating dog!",
		})
		return
	}

	c.JSON(http.StatusOK, io.MapToDogResponse(updatedDog, env.Storage))
}

// DogReUnited godoc
// @Summary Reunite dog with owner
// @Schemes
// @Description	Reunite dog with owner. Making him its only host and removing other hosts
// @Tags        dog
// @Accept      json
// @Produce     json
// @Param		dogID query string false "dog ID"
// @Param		ownerID query string false "dog owner ID"
// @Param		hostID query string false "dog host ID"
// @Success     200 {object} model.DogResponse
// @Failure		500 {object} map[string]string{error=string, message=string}
// @Router      /dog/found [patch]
func DogReUnited(c *gin.Context, env environment.Env) {
	q := c.Request.URL.Query()
	dogID, ownerID, hostID := q.Get("dogID"), q.Get("ownerID"), q.Get("hostID")
	dogIDInt, _ := strconv.Atoi(dogID)
	lfDogs := lostandfound.NewLostFoundDogs(env.DogPersister, env.UserPersister)

	dog, err := lfDogs.ReuniteDog(uint(dogIDInt), ownerID, hostID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "error updating lost dog status",
		})
		return
	}

	c.JSON(http.StatusOK, io.MapToDogResponse(dog, env.Storage))
}

// GetMissingDogsList godoc
// @Summary All missing dogs
// @Schemes
// @Description	List of all missing dogs
// @Tags        dog
// @Accept      json
// @Produce     json
// @Success     200 {object} []model.DogResponse
// @Router      /dog/missing [get]
func GetMissingDogsList(c *gin.Context, env environment.Env) {
	lfDogs := lostandfound.NewLostFoundDogs(env.DogPersister, env.UserPersister)
	dogList := lfDogs.GetMissingDogsList()
	dogRespList := io.MapToDogResponseList(dogList, env.Storage)

	c.JSON(http.StatusOK, dogRespList)
}

// ClaimFoundMissingDog godoc
// @Summary Claim that missing dog was found
// @Schemes
// @Description	Claim that missing dog is found
// @Tags        dog
// @Accept      json
// @Produce     json
// @Param		dogID query string false "dog ID"
// @Param		matchingDogs query array[string] false "possible matching dogs"
// @Success     200 string
// @Failure		500 {object} map[string]string{error=string, message=string}
// @Router      /dog/claim_found [patch]
func ClaimFoundMissingDog(c *gin.Context, env environment.Env) {
	q := c.Request.URL.Query()
	dogID, matchingDogsArr := q.Get("dogID"), q.Get("matchingDogs")
	dogIDInt, _ := strconv.Atoi(dogID)
	var matchingDogIDs []uint
	for _, matchingDogID := range strings.Split(matchingDogsArr, ",") {
		id, _ := strconv.Atoi(matchingDogID)
		matchingDogIDs = append(matchingDogIDs, uint(id))
	}
	lf := lostandfound.NewLostFoundDogs(env.DogPersister, env.UserPersister)
	err := lf.PossibleMatchingDogs(uint(dogIDInt), matchingDogIDs, users.NewUserManager(env.UserPersister), env.NotificationSender)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "error matching possible missing dogs",
		})
		return
	}
	c.String(http.StatusOK, "users notified!")
}
