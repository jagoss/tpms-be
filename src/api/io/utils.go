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
