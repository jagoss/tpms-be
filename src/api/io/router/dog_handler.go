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
// @Failure		422 {object} object{error=string, message=string}
// @Failure		500 {object} object{error=string, message=string}
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

// GetDog godoc
// @Summary Get dog given its ID
// @Schemes
// @Description Get dog given its ID
// @Tags        dog
// @Accept      json
// @Produce     json
// @Param		dog path string false  "dog ID"
// @Success     200 {object} model.DogResponse
// @Failure		400 {object} object{error=string, message=string}
// @Failure		500 {object} object{error=string, message=string}
// @Router      /dog/:id [get]
func GetDog(c *gin.Context, env environment.Env) {
	dogID, exists := c.Get("id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Variable missing",
			"message": "Missing dog ID",
		})
		return
	}
	dogManager := dogs.NewDogManager(env.DogPersister, env.Storage)

	dog, err := dogManager.Get(io.ParseToUint(dogID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Variable missing",
			"message": "Missing dog ID",
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
// @Failure		422 {object} object{error=string, message=string}
// @Failure		500 {object} object{error=string, message=string}
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
// @Failure		500 {object} object{error=string, message=string}
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
// @Summary Brings list of missing dogs
// @Schemes
// @Description	If no argument is given it returns all missing dogs. If user location and a search radius is sent, then it returns all missing dogs within that radius.
// @Tags        dog
// @Accept      json
// @Produce     json
// @Param		userLatitude query float64 false "user latitude"
// @Param		userLongitude query float64 false "user longitude"
// @Param		radius query float64 false "radio to look for dogs"
// @Success     200 {object} []model.DogResponse
// @Failure		400 {object} object{error=string, message=string}
// @Router      /dog/missing [get]
func GetMissingDogsList(c *gin.Context, env environment.Env) {
	q := c.Request.URL.Query()
	userLat, userLng, radius := q.Get("userLatitude"), q.Get("userLongitude"), q.Get("radius")
	if (userLat == "" || userLng == "" || radius == "") && (userLat != "" || userLng != "" || radius != "") {
		var missingArgs []string
		if userLat == "" {
			missingArgs = append(missingArgs, "UserLat")
		}
		if userLng == "" {
			missingArgs = append(missingArgs, "userLng")
		}
		if radius == "" {
			missingArgs = append(missingArgs, "radius")
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Missing argument",
			"message": fmt.Sprintf("There are missing arguments: %s", strings.Join(missingArgs, ",")),
		})
		return
	}

	lfDogs := lostandfound.NewLostFoundDogs(env.DogPersister, env.UserPersister)
	if userLat == "" && userLng == "" && radius == "" {
		dogList := lfDogs.GetAllMissingDogsList()
		dogRespList := io.MapToDogResponseList(dogList, env.Storage)
		c.JSON(http.StatusOK, dogRespList)
		return
	}
	userLatF, _ := strconv.ParseFloat(userLat, 64)
	userLngF, _ := strconv.ParseFloat(userLng, 64)
	radiusF, _ := strconv.ParseFloat(radius, 64)
	dogList := lfDogs.GetMissingDogsInRadius(userLatF, userLngF, radiusF)
	dogRespList := io.MapToDogResponseList(dogList, env.Storage)

	c.JSON(http.StatusOK, dogRespList)
}
