package router

import (
	"be-tpms/src/api/environment"
	"be-tpms/src/api/io"
	"be-tpms/src/api/io/fileio"
	"be-tpms/src/api/usecases/dogs"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
)

type DogHandler struct {
}

func RegisterNewDog(c *gin.Context, env environment.Env) {
	jsonBody, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Printf("error reading request body: %v", err)
		c.String(http.StatusUnprocessableEntity, "error reading request body!")
	}
	reqDog, err := io.DeserializeDog(jsonBody)
	if err != nil {
		log.Printf("error unmarshalling dog body: %v", err)
		c.String(http.StatusUnprocessableEntity, "error reading dog body!")
	}
	dogManager := dogs.NewDogManager(env.DogPersister)
	newDog, img := io.MapFromDogRequest(reqDog)
	imgURL, err := fileio.SaveImgs(img)
	if err != nil {
		log.Printf("%v", err)
		c.String(http.StatusInternalServerError, fmt.Sprintf("error saving img: %v", err))
		return
	}
	newDog.ImgUrl = imgURL
	dog, err := dogManager.Register(newDog) //agregar img
	if err != nil {
		log.Printf("%v", err)
		c.String(http.StatusInternalServerError, fmt.Sprintf("error inserting new dog: %v", err))
	}

	c.JSON(http.StatusOK, dog)
}

func UpdateDog(c *gin.Context, env environment.Env) {
	jsonBody, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Printf("error reading request body: %v", err)
		c.String(http.StatusUnprocessableEntity, "error reading request body!")
	}
	dogReq, err := io.DeserializeDog(jsonBody)
	if err != nil {
		log.Printf("error unmarshalling dog body: %v", err)
		c.String(http.StatusUnprocessableEntity, "error reading dog body!")
	}
	dogManager := dogs.NewDogManager(env.DogPersister)
	dog, _ := io.MapFromDogRequest(dogReq)
	updatedDog, err := dogManager.Modify(dog)
	if err != nil {
		log.Printf("error updating dog with ID %d: %v ", dog.ID, err)
		c.String(http.StatusInternalServerError, "error updating dog")
	}
	c.JSON(http.StatusOK, updatedDog)
}

func ReportALostDog(c *gin.Context, env environment.Env) {
	jsonBody, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Printf("error reading request body: %v", err)
		c.String(http.StatusUnprocessableEntity, "error reading request body!")
	}
	lostDogReq, err := io.DeserializeDog(jsonBody)
	if err != nil {
		log.Printf("error unmarshalling dog body: %v", err)
		c.String(http.StatusUnprocessableEntity, "error reading dog body!")
	}
	dogManager := dogs.NewDogManager(env.DogPersister)
	lostDog, _ := io.MapFromDogRequest(lostDogReq)
	dog, err := dogManager.ReportLostDog(lostDog)
	if err != nil {
		log.Printf("%v", err)
		c.String(http.StatusInternalServerError, "error reporting lost dog")
	}

	c.JSON(http.StatusOK, dog)
}
