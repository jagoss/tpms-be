package io

import (
	"be-tpms/src/api/domain/model"
	"be-tpms/src/api/usecases/interfaces"
	"encoding/json"
	"fmt"
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

func MapFromDogRequest(reqDog *model.DogRequest) (*model.Dog, []string) {
	dog := &model.Dog{
		Name:       reqDog.Name,
		Breed:      model.ParseBreed(reqDog.Breed),
		Age:        model.ParseAge(reqDog.Age),
		Size:       model.ParseSize(reqDog.Size),
		CoatColor:  model.ParseCoatColor(reqDog.CoatColor),
		CoatLength: model.ParseCoatLength(reqDog.CoatLength),
		IsLost:     reqDog.IsLost,
		Owner:      &model.User{ID: reqDog.Owner},
		Host:       &model.User{ID: reqDog.Host},
		Latitude:   reqDog.Latitude,
		Longitude:  reqDog.Longitude,
		ImgUrl:     reqDog.ImgUrl,
	}
	unitID, err := strconv.ParseUint(reqDog.ID, 10, 64)
	if err != nil {
		return nil, nil
	}
	dog.ID = uint(unitID)

	return dog, reqDog.Imgs
}

func MapToDogResponse(dog *model.Dog, bucket interfaces.Storage) *model.DogResponse {
	firstImg := strings.Split(dog.ImgUrl, ";")[0]
	imgArray, _ := bucket.GetImgs(firstImg)
	return &model.DogResponse{
		ID:         strconv.Itoa(int(dog.ID)),
		Name:       dog.Name,
		Breed:      dog.Breed.String(),
		Age:        dog.Age.String(),
		Size:       dog.Size.String(),
		CoatColor:  dog.CoatColor.String(),
		CoatLength: dog.CoatLength.String(),
		IsLost:     dog.IsLost,
		Owner:      dog.Owner.ID,
		Host:       dog.Host.ID,
		Latitude:   dog.Latitude,
		Longitude:  dog.Longitude,
		ImgsUrl:    dog.ImgUrl,
		ProfileImg: imgArray[0],
	}
}

func MapToDogResponseList(dogs []model.Dog, bucket interfaces.Storage) []model.DogResponse {
	if len(dogs) == 0 {
		return nil
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
