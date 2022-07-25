package router

import (
	"be-tpms/src/api/environment"
	"be-tpms/src/api/io"
	"be-tpms/src/api/usecases/dogs"
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
	newDog, err := io.DeserializeDog(jsonBody)
	if err != nil {
		log.Printf("error unmarshalling dog body: %v", err)
		c.String(http.StatusUnprocessableEntity, "error reading dog body!")
	}
	dogManager := dogs.NewDogManager(env.DogPersister)
	dog, err := dogManager.Register(newDog)
	if err != nil {
		log.Printf("%v", err)
		c.String(http.StatusInternalServerError, "error inserting new dog")
	}

	c.JSON(http.StatusOK, dog)
}

func UpdateDog(c *gin.Context, env environment.Env) {
	jsonBody, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Printf("error reading request body: %v", err)
		c.String(http.StatusUnprocessableEntity, "error reading request body!")
	}
	dog, err := io.DeserializeDog(jsonBody)
	if err != nil {
		log.Printf("error unmarshalling dog body: %v", err)
		c.String(http.StatusUnprocessableEntity, "error reading dog body!")
	}
	dogManager := dogs.NewDogManager(env.DogPersister)
	updatedDog, err := dogManager.Modify(dog)
	if err != nil {
		log.Printf("error updating dog with ID %s: %v ", dog.ID, err)
		c.String(http.StatusInternalServerError, "error updating dog")
	}
	c.JSON(http.StatusOK, updatedDog)
}
