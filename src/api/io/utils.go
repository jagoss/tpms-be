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

func DeserializeDog(input []byte) (*model.Dog, error) {
	var dog model.Dog
	err := json.Unmarshal(input, &dog)
	if err != nil {
		return nil, err
	}
	return &dog, nil
}
