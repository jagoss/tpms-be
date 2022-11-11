package io

import (
	"be-tpms/src/api/domain/model"
	"be-tpms/src/api/usecases/interfaces"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
)

const (
	TITLE = "title"
	BODY  = "body"
)

func DeserializeUser(input []byte) (*model.User, error) {
	var user model.User
	err := json.Unmarshal(input, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func DeserializeDog(input []byte) (*model.DogRequest, error) {
	var dog model.DogRequest
	err := json.Unmarshal(input, &dog)
	if err != nil {
		return nil, err
	}
	return &dog, nil
}

func MapFromDogRequest(reqDog *model.DogRequest) (*model.Dog, []string) {
	dog := &model.Dog{
		Name:       reqDog.Name,
		Breed:      model.ParseBreed(reqDog.Breed),
		Age:        model.ParseAge(reqDog.Age),
		Size:       model.ParseSize(reqDog.Size),
		CoatColor:  model.ParseCoatColor(reqDog.CoatColor),
		CoatLength: model.ParseCoatLength(reqDog.CoatLength),
		TailLength: model.ParseTailLength(reqDog.TailLength),
		IsLost:     reqDog.IsLost,
		Latitude:   reqDog.Latitude,
		Longitude:  reqDog.Longitude,
		ImgUrl:     reqDog.ImgUrl,
	}
	if reqDog.Owner != "" {
		dog.Owner = &model.User{ID: reqDog.Owner}
	}
	if reqDog.Host != "" {
		dog.Host = &model.User{ID: reqDog.Host}
	}

	if reqDog.ID != "" {
		unitID, err := strconv.ParseInt(reqDog.ID, 10, 64)
		if err != nil {
			log.Printf("error parsing uint: %s", err.Error())
			return nil, nil
		}
		dog.ID = unitID
	}

	return dog, reqDog.Imgs
}

func MapToDogResponse(dog *model.Dog, bucket interfaces.Storage) *model.DogResponse {
	response := &model.DogResponse{
		ID:         strconv.Itoa(int(dog.ID)),
		Name:       dog.Name,
		Breed:      dog.Breed.String(),
		Age:        dog.Age.String(),
		Size:       dog.Size.String(),
		CoatColor:  dog.CoatColor.String(),
		CoatLength: dog.CoatLength.String(),
		TailLength: dog.TailLength.String(),
		IsLost:     dog.IsLost,
		Latitude:   dog.Latitude,
		Longitude:  dog.Longitude,
		ImgsUrl:    dog.ImgUrl,
	}

	if dog.Owner != nil {
		response.Owner = dog.Owner.ID
	}
	if dog.Host != nil {
		response.Host = dog.Host.ID
	}

	firstImg := strings.Split(dog.ImgUrl, ";")[0]
	imgArray, _ := bucket.GetImgs(firstImg)
	if len(imgArray) != 0 {
		response.ProfileImg = imgArray[0]
	}

	return response
}

func MapToDogResponseList(dogs []model.Dog, bucket interfaces.Storage) []model.DogResponse {
	if len(dogs) == 0 {
		return make([]model.DogResponse, 0)
	}
	var dogsResp []model.DogResponse
	for _, dog := range dogs {
		dogsResp = append(dogsResp, *MapToDogResponse(&dog, bucket))
	}
	return dogsResp
}

func ParseToUint(val any) uint {
	stringVal := fmt.Sprintf("%s", val)
	res, _ := strconv.ParseUint(stringVal, 10, 64)
	return uint(res)
}

func ToStringList(values []uint) []string {
	var resultList []string
	for _, val := range values {
		resultList = append(resultList, fmt.Sprintf("%d", val))
	}
	if len(resultList) == 0 {
		return make([]string, 0)
	}
	return resultList
}
