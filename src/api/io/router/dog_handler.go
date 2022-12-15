package router

import (
	"be-tpms/src/api/domain/model"
	"be-tpms/src/api/environment"
	"be-tpms/src/api/io"
	"be-tpms/src/api/usecases/cvmodel"
	"be-tpms/src/api/usecases/dogs"
	"be-tpms/src/api/usecases/lostandfound"
	"be-tpms/src/api/usecases/messaging"
	"be-tpms/src/api/usecases/users"
	"encoding/json"
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
// @Failure		400 {object} object{error=string,message=string}
// @Failure		401 {object} object{error=string,message=string}
// @Failure		422 {object} object{error=string,message=string}
// @Failure		500 {object} object{error=string,message=string}
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
		log.Printf("error mapping dog request when parsing uint")
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   "error parsing uint",
			"message": "error parsing dogRequest to dog!",
		})
		return
	}
	log.Printf("registring dog")
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
	log.Printf("go to send notifications")
	if dog.IsLost {
		notificationSender := messaging.NewMessageSender(env.NotificationSender, env.UserPersister)
		if err = notificationSender.SendToEnabledUsers(dog); err != nil {
			log.Printf("error notifying users")
		}
	}
	log.Printf("dog register, calculating embdding")
	predictionService := cvmodel.NewPrediction(env.DogPersister, env.CVModelRestClient, env.Storage)
	if err = predictionService.CalculateEmbedding(uint(dog.ID)); err != nil {
		log.Printf("%v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": fmt.Sprintf("error calculating new dog %d vector: %v", dog.ID, err),
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
// @Failure		400 {object} object{error=string,message=string}
// @Failure		401 {object} object{error=string,message=string}
// @Failure		500 {object} object{error=string,message=string}
// @Router      /dog/:id [get]
func GetDog(c *gin.Context, env environment.Env) {
	dogID := c.Param("id")
	if dogID == "" {
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
			"error":   err.Error(),
			"message": "error getting dog",
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
// @Failure		400 {object} object{error=string,message=string}
// @Failure		401 {object} object{error=string,message=string}
// @Failure		422 {object} object{error=string,message=string}
// @Failure		500 {object} object{error=string,message=string}
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
// @Param		matchingDog query string false "matching dog ID"
// @Success     200 {object} model.DogResponse
// @Failure		400 {object} object{error=string,message=string}
// @Failure		401 {object} object{error=string,message=string}
// @Failure		500 {object} object{error=string,message=string}
// @Router      /dog/found [put]
func DogReUnited(c *gin.Context, env environment.Env) {
	q := c.Request.URL.Query()
	dogID, possibleDogID := q.Get("dogId"), q.Get("possibleDogId")
	dogIDInt, _ := strconv.Atoi(dogID)
	possibleDogIDInt, _ := strconv.Atoi(possibleDogID)

	lfDogs := lostandfound.NewLostFoundDogs(env.DogPersister, env.UserPersister, env.PossibleMatchPersister)

	dog, err := lfDogs.ReuniteDog(uint(dogIDInt), uint(possibleDogIDInt), env.NotificationSender)
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
// @Failure		400 {object} object{error=string,message=string}
// @Failure		401 {object} object{error=string,message=string}
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

	lfDogs := lostandfound.NewLostFoundDogs(env.DogPersister, env.UserPersister, env.PossibleMatchPersister)
	if userLat == "" && userLng == "" && radius == "" {
		dogList, err := lfDogs.GetAllMissingDogsList()
		if err != nil {
			log.Printf(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   err.Error(),
				"message": "error getting all lost dogs",
			},
			)
			return
		}
		dogRespList := io.MapToDogResponseList(dogList, env.Storage)
		c.JSON(http.StatusOK, dogRespList)
		return
	}
	userLatF, _ := strconv.ParseFloat(userLat, 64)
	userLngF, _ := strconv.ParseFloat(userLng, 64)
	radiusF, _ := strconv.ParseFloat(radius, 64)
	dogList, err := lfDogs.GetMissingDogsInRadius(userLatF, userLngF, radiusF)
	if err != nil {
		log.Printf(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": fmt.Sprintf("error getting lost dogs in radius of %2.f", radiusF),
		},
		)
		return
	}

	dogRespList := io.MapToDogResponseList(dogList, env.Storage)

	c.JSON(http.StatusOK, dogRespList)
}

// PossibleMatch godoc
// @Summary    Mark dog as possible dog
// @Schemes
// @Description	Mark dog as possible match and notify host of that dog if exists.
// @Tags        dog
// @Accept      json
// @Produce     json
// @Param		dogID body string false "dog ID"
// @Param		possibleDogs body []string false "possible matching dogs"
// @Failure		400 {object} object{error=string,message=string}
// @Failure		401 {object} object{error=string,message=string}
// @Failure		500 {object} object{error=string,message=string}
// @Success     200 {object} object{message=string}
// @Router      /dog/possible [post]
func PossibleMatch(c *gin.Context, env environment.Env) {
	jsonBody, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Printf("error reading request body: %v", err)
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   err.Error(),
			"message": "error reading request body!",
		})
		return
	}
	var body map[string]interface{}
	err = json.Unmarshal(jsonBody, &body)
	if err != nil {
		log.Printf("unmarshalling error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "error unmarshalling request body!",
		})
		return
	}

	log.Printf("[PossibleMatch] request body: %v", &body)

	if body["dogId"] == nil {
		log.Printf("status code 400: missing dogId")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "missing dogId",
			"message": "missing key value",
		})
		return
	}
	if body["possibleDogs"] == nil {
		log.Printf("status code 400: missing possibleDogs")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "missing possibleDogs",
			"message": "missing key value",
		})
		return
	}

	dogID, possibleDogIDs := fmt.Sprintf("%v", body["dogId"]), io.ToArray(body["possibleDogs"])
	dogIDInt, _ := strconv.Atoi(dogID)

	log.Printf("possibleDogs: %v", &possibleDogIDs)
	var matchingDogIDs []uint
	for _, matchingDogID := range possibleDogIDs {
		id, _ := strconv.Atoi(matchingDogID)
		matchingDogIDs = append(matchingDogIDs, uint(id))
	}

	lf := lostandfound.NewLostFoundDogs(env.DogPersister, env.UserPersister, env.PossibleMatchPersister)
	err = lf.PossibleMatchingDogs(uint(dogIDInt), matchingDogIDs, env.NotificationSender)
	if err != nil {
		log.Printf(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "error matching possible missing dogs",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "users notified!"})
}

// AckPossibleDog godoc
// @Summary Acknowledge possible dog as the dog
// @Schemes
// @Description	Acknowledge that possible dog is the missing dog
// @Tags        dog
// @Accept      json
// @Produce     json
// @Param		dogId body string false "dog ID"
// @Param		possibleDogId body string false "possible dog ID"
// @Failure		400 {object} object{error=string,message=string}
// @Failure		401 {object} object{error=string,message=string}
// @Failure		500 {object} object{error=string,message=string}
// @Success     200 {object} object{message=string}
// @Router      /dog/possible [put]
func AckPossibleDog(c *gin.Context, env environment.Env) {
	jsonBody, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Printf("error reading request body: %v", err)
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   err.Error(),
			"message": "error reading request body!",
		})
		return
	}
	var body map[string]interface{}
	err = json.Unmarshal(jsonBody, &body)
	if err != nil {
		log.Printf("unmarshalling error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "error unmarshalling request body!",
		})
		return
	}
	if body["dogId"] == nil {
		log.Printf("status code 400: missing dogId")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "missing dogId",
			"message": "missing key value",
		})
		return
	}
	if body["possibleDogId"] == nil {
		log.Printf("status code 400: missing possibleDogId")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "missing possibleDogId",
			"message": "missing key value",
		})
		return
	}
	dogID, possibleDogID := fmt.Sprintf("%v", body["dogId"]), fmt.Sprintf("%v", body["possibleDogId"])
	dogIDInt, _ := strconv.Atoi(dogID)
	possibleDogIDInt, _ := strconv.Atoi(possibleDogID)

	lf := lostandfound.NewLostFoundDogs(env.DogPersister, env.UserPersister, env.PossibleMatchPersister)
	if err := lf.AcknowledgePossibleDog(uint(dogIDInt), uint(possibleDogIDInt), env.NotificationSender); err != nil {
		log.Printf(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "error matching possible missing dogs",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "users notified!"})
}

// RejectPossibleDog godoc
// @Summary Reject possible dog as dog
// @Schemes
// @Description	Reject that possible dog is missing dog. Delete register from table. Then notify user
// @Tags        dog
// @Accept      json
// @Produce     json
// @Param		dogID query string false "dog ID"
// @Param		possibleDogID query string false "possible dog ID"
// @Failure		400 {object} object{error=string,message=string}
// @Failure		401 {object} object{error=string,message=string}
// @Failure		500 {object} object{error=string,message=string}
// @Success     200 {object} object{message=string}
// @Router      /dog/possible [delete]
func RejectPossibleDog(c *gin.Context, env environment.Env) {
	jsonBody, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Printf("error reading request body: %v", err)
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   err.Error(),
			"message": "error reading request body!",
		})
		return
	}
	var body map[string]interface{}
	err = json.Unmarshal(jsonBody, &body)
	if err != nil {
		log.Printf("unmarshalling error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "error unmarshalling request body!",
		})
		return
	}
	if body["dogId"] == nil {
		log.Printf("status code 400: missing dogId")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "missing dogId",
			"message": "missing key value",
		})
		return
	}
	if body["possibleDogId"] == nil {
		log.Printf("status code 400: missing possibleDogId")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "missing possibleDogId",
			"message": "missing key value",
		})
		return
	}
	dogID, possibleDogID := fmt.Sprintf("%v", body["dogId"]), fmt.Sprintf("%v", body["possibleDogId"])
	dogIDInt, _ := strconv.Atoi(dogID)
	possibleDogIDInt, _ := strconv.Atoi(possibleDogID)

	lf := lostandfound.NewLostFoundDogs(env.DogPersister, env.UserPersister, env.PossibleMatchPersister)
	err = lf.RejectPossibleDog(uint(dogIDInt), uint(possibleDogIDInt), env.NotificationSender)
	if err != nil {
		log.Printf("error: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "error matching possible missing dogs",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "users notified!"})
}

// DeleteDog godoc
// @Summary Delete dog given its ID
// @Schemes
// @Description Delete dog given its ID
// @Tags        dog
// @Accept      json
// @Produce     json
// @Param		dog path string false  "dog ID"
// @Success     200 {object} object{deleted=bool}
// @Failure		400 {object} object{error=string,message=string}
// @Failure		500 {object} object{error=string,message=string}
// @Router      /dog/:id [delete]
func DeleteDog(c *gin.Context, env environment.Env) {
	dogID := c.Param("id")
	if dogID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Variable missing",
			"message": "Missing dog ID",
		})
		return
	}

	dogManager := dogs.NewDogManager(env.DogPersister, env.Storage)
	deleted, err := dogManager.Delete(io.ParseToUint(dogID), env.PossibleMatchPersister)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "error deleting dog",
		})
		return
	}
	if !deleted {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "could not delete dog",
			"message": "error deleting dog",
		})
	}
	c.JSON(http.StatusOK, gin.H{"delted": deleted})
}

// GetPossibleMatchingDogs godoc
// @Summary Get possible matching dogs given dog id and ack status
// @Schemes
// @Description Given one dog ID return possible matching dogs and ack status of confirmation
// @Tags        dog
// @Accept      json
// @Produce     json
// @Param		dog path string false  "dog ID"
// @Param		acks query []string	true "matching confirmation status"
// @Success     200 {object} []model.PossibleMatch
// @Failure		400 {object} object{error=string,message=string}
// @Failure		401 {object} object{error=string,message=string}
// @Failure		500 {object} object{error=string,message=string}
// @Router      /dog/:id/possible [get]
func GetPossibleMatchingDogs(c *gin.Context, env environment.Env) {
	id := c.Param("id")
	acksStrings, exists := c.GetQueryArray("acks")

	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Variable missing",
			"message": "Missing acks",
		})
		return
	}

	var acks []model.Ack
	for _, ack := range acksStrings {
		acks = append(acks, model.ParseAck(ack))
	}

	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Variable missing",
			"message": "Missing dog ID",
		})
		return
	}
	lfDogs := lostandfound.NewLostFoundDogs(env.DogPersister, env.UserPersister, env.PossibleMatchPersister)
	dogID, _ := strconv.ParseUint(id, 10, 64)
	possibleMatchingDogs, err := lfDogs.GetPossibleMatchingDogs(uint(dogID), acks)

	if err != nil {
		log.Printf(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": fmt.Sprintf("error getting lost dogs in radius of %s", id),
		})
		return
	}
	log.Printf("[GetPossibleMatchingDogs] result list: %v", possibleMatchingDogs)
	c.JSON(http.StatusOK, io.PossibleMatchListDto(possibleMatchingDogs))
}

// GetSimilarDogPrediction godoc
// @Summary Get similar dogs from CV models
// @Schemes
// @Description Given one dog ID return similar based on results of CV prediction model
// @Tags        dog
// @Accept      json
// @Produce     json
// @Param		dog path string false  "dog ID"
// @Success     200 {object} []model.DogResponse
// @Failure		400 {object} object{error=string,message=string}
// @Failure		401 {object} object{error=string,message=string}
// @Failure		500 {object} object{error=string,message=string}
// @Router      /dog/:id/prediction [get]
func GetSimilarDogPrediction(c *gin.Context, env environment.Env) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Variable missing",
			"message": "Missing dog ID",
		})
		return
	}
	predictionService := cvmodel.NewPrediction(env.DogPersister, env.CVModelRestClient, env.Storage)
	dogID, _ := strconv.ParseUint(id, 10, 64)
	resultList, err := predictionService.FindMatches(uint(dogID))
	if err != nil {
		log.Printf(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": fmt.Sprintf("error getting similar dogs from CV model for dog %s", id),
		})
		return
	}

	c.JSON(http.StatusOK, io.MapToDogResponseList(resultList, env.Storage))
}

// ReportDogAsMissing godoc
// @Summary Update dog status as missing
// @Schemes
// @Description Update given dog status as missing
// @Tags        dog
// @Accept      json
// @Produce     json
// @Param		dog path string false  "dog ID"
// @Param		userLatitude query float64 false "user lat"
// @Param		userLongitude query float64 false "user lng"
// @Success     200 {object} []model.DogResponse
// @Failure		400 {object} object{error=string,message=string}
// @Failure		401 {object} object{error=string,message=string}
// @Failure		500 {object} object{error=string,message=string}
// @Router      /dog/:id [patch]
func ReportDogAsMissing(c *gin.Context, env environment.Env) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Variable missing",
			"message": "Missing dog ID",
		})
		return
	}

	lat, lng := c.Query("lat"), c.Query("lng")
	if lat == "" || lng == "" {
		var missingArgs []string
		if lat == "" {
			missingArgs = append(missingArgs, "lat")
		}
		if lng == "" {
			missingArgs = append(missingArgs, "lng")
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Missing argument",
			"message": fmt.Sprintf("There are missing arguments: %s", strings.Join(missingArgs, ",")),
		})
		return
	}

	dogManager := dogs.NewDogManager(env.DogPersister, env.Storage)
	idUint, _ := strconv.ParseUint(id, 10, 64)
	latFloat, _ := strconv.ParseFloat(lat, 64)
	lngFloat, _ := strconv.ParseFloat(lng, 64)
	dog, err := dogManager.SetDogAsLost(uint(idUint), latFloat, lngFloat)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": fmt.Sprintf("could not update dog %s", id),
		})
		return
	}

	notificationSender := messaging.NewMessageSender(env.NotificationSender, env.UserPersister)
	if err = notificationSender.SendToEnabledUsers(dog); err != nil {
		log.Printf("error notifying users")
	}

	predictionService := cvmodel.NewPrediction(env.DogPersister, env.CVModelRestClient, env.Storage)
	dogID, _ := strconv.ParseUint(id, 10, 64)
	if err := predictionService.CalculateEmbedding(uint(dogID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": fmt.Sprintf("could not update dog %s", id),
		})
		return
	}

	c.JSON(http.StatusOK, io.MapToDogResponse(dog, env.Storage))
}

func GenerarteEmbedding(c *gin.Context, env environment.Env) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Variable missing",
			"message": "Missing dog ID",
		})
		return
	}
	predictionService := cvmodel.NewPrediction(env.DogPersister, env.CVModelRestClient, env.Storage)
	dogID, _ := strconv.ParseUint(id, 10, 64)
	if err := predictionService.CalculateEmbedding(uint(dogID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": fmt.Sprintf("could not update dog %s", id),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Embedding calculated for dog %s", id)})
}

// RegisterNewScrapperDog godoc
// @Summary Register new dog from scrapper search
// @Schemes
// @Description Register new dog from scrapper search
// @Tags        dog
// @Accept      json
// @Produce     json
// @Param		dog body model.DogRequest false  "dog"
// @Success     200 {object} model.DogResponse
// @Failure		400 {object} object{error=string,message=string}
// @Failure		401 {object} object{error=string,message=string}
// @Failure		422 {object} object{error=string,message=string}
// @Failure		500 {object} object{error=string,message=string}
// @Router      /dog/scrapper [post]
func RegisterNewScrapperDog(c *gin.Context, env environment.Env) {
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
		log.Printf("error mapping dog request when parsing uint")
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   "error parsing uint",
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
