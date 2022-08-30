package router

import (
	"be-tpms/src/api/environment"
	"be-tpms/src/api/io"
	"be-tpms/src/api/usecases/dogs"
	"be-tpms/src/api/usecases/lostandfound"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

type DogHandler struct {
}

func RegisterNewDog(c *gin.Context, env environment.Env) {
	jsonBody, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Printf("error reading request body: %v", err)
		c.String(http.StatusUnprocessableEntity, "error reading request body!")
		return
	}
	reqDog, err := io.DeserializeDog(jsonBody)
	if err != nil {
		log.Printf("error unmarshalling dog body: %v", err)
		c.String(http.StatusUnprocessableEntity, "error reading dog body!")
		return
	}
	dogManager := dogs.NewDogManager(env.DogPersister, env.Storage)
	bucket := env.Storage
	imgURL, err := bucket.SaveImgs(reqDog.Imgs)
	if err != nil {
		log.Printf("%v", err)
		c.String(http.StatusInternalServerError, fmt.Sprintf("error saving img: %v", err))
		return
	}
	reqDog.Dog.ImgUrl = imgURL
	dog, err := dogManager.Register(&reqDog.Dog) //agregar img
	if err != nil {
		log.Printf("%v", err)
		c.String(http.StatusInternalServerError, fmt.Sprintf("error inserting new dog: %v", err))
		return
	}

	c.JSON(http.StatusOK, dog)
}

func UpdateDog(c *gin.Context, env environment.Env) {
	jsonBody, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Printf("error reading request body: %v", err)
		c.String(http.StatusUnprocessableEntity, "error reading request body!")
		return
	}
	dogReq, err := io.DeserializeDog(jsonBody)
	if err != nil {
		log.Printf("error unmarshalling dog body: %v", err)
		c.String(http.StatusUnprocessableEntity, "error reading dog body!")
		return
	}
	dogManager := dogs.NewDogManager(env.DogPersister, env.Storage)
	updatedDog, err := dogManager.Modify(&dogReq.Dog, dogReq.Imgs)
	if err != nil {
		log.Printf("error updating dog with ID %d: %v ", dogReq.Dog.ID, err)
		c.String(http.StatusInternalServerError, "error updating dog")
		return
	}
	c.JSON(http.StatusOK, updatedDog)
}

func DogReUnited(c *gin.Context, env environment.Env) {
	q := c.Request.URL.Query()
	dogID, ownerID, hostID := q.Get("dogID"), q.Get("ownerID"), q.Get("hostID")
	dogIDInt, _ := strconv.Atoi(dogID)
	lfDogs := lostandfound.NewLostFoundDogs(env.DogPersister, env.UserPersister)
	dog, err := lfDogs.ReuniteDog(uint(dogIDInt), ownerID, hostID)
	if err != nil {
		c.String(http.StatusInternalServerError, "error updating lost dog status")
		return
	}
	c.JSON(http.StatusOK, dog)
}

func GetMissingDogsList(c *gin.Context, env environment.Env) {
	lfDogs := lostandfound.NewLostFoundDogs(env.DogPersister, env.UserPersister)
	dogList := lfDogs.GetMissingDogsList()
	dogRespList := io.MapToDogResponse(dogList, env.Storage)
	c.JSON(http.StatusOK, dogRespList)
}
