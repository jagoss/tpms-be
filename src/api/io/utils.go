package io

import (
	"be-tpms/src/api/domain/model"
	"encoding/json"
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
