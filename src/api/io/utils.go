package io

import (
	"be-tpms/src/api/domain/model"
	"be-tpms/src/api/usecases/interfaces"
	"encoding/json"
	"strconv"
	"strings"
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

func MapFromDogRequest(reqDog *model.DogRequest) (*model.Dog, [][]byte) {
	return &model.Dog{
		Name:      "",
		Breed:     0,
		Age:       0,
		Size:      0,
		Owner:     nil,
		Host:      nil,
		Latitude:  0,
		Longitude: 0,
		ImgUrl:    "",
	}, reqDog.Img
}

func MapToDogResponse(dogs []model.Dog, bucket interfaces.Storage) []model.DogResponse {
	var dogsResp []model.DogResponse
	for _, dog := range dogs {
		firstImg := strings.Split(dog.ImgUrl, ";")[0]
		imgArray, _ := bucket.GetImgs(firstImg)
		dogsResp = append(dogsResp, model.DogResponse{ID: strconv.Itoa(int(dog.ID)), Dog: dog, Img: imgArray[0]})
	}
	return dogsResp
}
