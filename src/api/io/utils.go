package io

import (
	"be-tpms/src/api/domain/model"
	"be-tpms/src/api/usecases/interfaces"
	"encoding/json"
	"fmt"
	"log"
	"reflect"
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

func DeserializePosts(input []byte) ([]model.PostRequest, error) {
	var posts []model.PostRequest
	err := json.Unmarshal(input, &posts)
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func MapFromDogRequest(reqDog *model.DogRequest) (*model.Dog, []string) {
	dog := &model.Dog{
		Name:           reqDog.Name,
		Breed:          model.ParseBreed(reqDog.Breed),
		Age:            model.ParseAge(reqDog.Age),
		Size:           model.ParseSize(reqDog.Size),
		CoatColor:      model.ParseCoatColor(reqDog.CoatColor),
		CoatLength:     model.ParseCoatLength(reqDog.CoatLength),
		TailLength:     model.ParseTailLength(reqDog.TailLength),
		Ear:            model.ParseEar(reqDog.Ear),
		AdditionalInfo: reqDog.AdditionalInfo,
		IsLost:         reqDog.IsLost,
		Latitude:       reqDog.Latitude,
		Longitude:      reqDog.Longitude,
		ImgUrl:         reqDog.ImgUrl,
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
		ID:             strconv.Itoa(int(dog.ID)),
		Name:           dog.Name,
		Breed:          dog.Breed.String(),
		Age:            dog.Age.String(),
		Size:           dog.Size.String(),
		CoatColor:      dog.CoatColor.String(),
		CoatLength:     dog.CoatLength.String(),
		TailLength:     dog.TailLength.String(),
		Ear:            dog.Ear.String(),
		AdditionalInfo: dog.AdditionalInfo,
		IsLost:         dog.IsLost,
		Latitude:       dog.Latitude,
		Longitude:      dog.Longitude,
		ImgsUrl:        dog.ImgUrl,
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

func ParseToUintList(list []int64) []uint {
	var resultList []uint
	for _, e := range list {
		resultList = append(resultList, ParseToUint(e))
	}
	return resultList
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

func PossibleMatchListDto(matches []model.PossibleMatch) []model.PossibleMatchDto {
	var resultList []model.PossibleMatchDto
	for _, pm := range matches {
		resultList = append(resultList, PossibleMatchToDto(pm))
	}
	log.Printf("possible match dto: %v", resultList)
	if len(resultList) == 0 {
		return make([]model.PossibleMatchDto, 0)
	}
	return resultList
}

func PossibleMatchToDto(match model.PossibleMatch) model.PossibleMatchDto {
	return model.PossibleMatchDto{
		DogID:         strconv.Itoa(int(match.DogID)),
		PossibleDogID: strconv.Itoa(int(match.PossibleDogID)),
		Ack:           match.Ack.String(),
	}
}

func ToArray(t interface{}) []string {
	var result []string
	log.Printf("[io.ToArray] type: %v", reflect.TypeOf(t).Kind())
	switch reflect.TypeOf(t).Kind() {
	case reflect.Slice:
		s := t.([]interface{})
		for i := 0; i < len(s); i++ {
			result = append(result, fmt.Sprintf("%s", s[i]))
		}
		break
	case reflect.String:
		s := t.(string)
		s = strings.Replace(strings.Replace(s, "[", "", 1), "]", "", 1)
		result = strings.Split(s, ",")
	}
	log.Printf("[io.ToArray] result: %s", result)
	return result
}

func MapFromPostRequest(postReq *model.PostRequest) *model.Post {
	return &model.Post{
		Url:      postReq.Url,
		Title:    postReq.Title,
		Location: postReq.Location,
	}
}

func MapFromPostRequestList(postReqList []model.PostRequest) []model.Post {
	if postReqList == nil || len(postReqList) == 0 {
		return make([]model.Post, 0)
	}
	var resultList []model.Post
	for _, p := range postReqList {
		resultList = append(resultList, *MapFromPostRequest(&p))
	}
	return resultList
}

func MapToPostResponse(post *model.Post, name string, img string, bucket interfaces.Storage) *model.PostResponse {
	p := model.PostResponse{
		Id:       strconv.FormatInt(post.Id, 10),
		DogId:    strconv.FormatInt(post.DogId, 10),
		DogName:  name,
		Url:      post.Url,
		Title:    post.Title,
		Location: post.Location,
	}
	imgArray, _ := bucket.GetImgs(img)
	p.Image = imgArray[0]
	return &p
}

func MapToPostResponseList(posts []model.Post, dogList []model.Dog, bucket interfaces.Storage) []model.PostResponse {
	if len(posts) == 0 {
		return make([]model.PostResponse, 0)
	}
	var postsResp []model.PostResponse
	for i, post := range posts {
		postsResp = append(postsResp, *MapToPostResponse(&post, dogList[i].Name, dogList[i].ImgUrl, bucket))
	}
	return postsResp
}
