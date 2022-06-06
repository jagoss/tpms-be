package restclient

import (
	"be-tpms/src/api/domain/model"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
)

const (
	baseURL = "localhost:8081"
)

type CVModelClient struct {
	rc *resty.Client
}

func NewCVModelRestClient(client *resty.Client) *CVModelClient {
	return &CVModelClient{rc: client}
}

func (client *CVModelClient) SearchDog() (*model.DogResponse, error) {
	response, err := client.rc.R().
		SetHeader("Accept", "application/json").
		Get(baseURL + "/dog")
	if err != nil {
		return nil, err
	}
	var dog model.DogResponse
	err = json.Unmarshal(response.Body(), &dog)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling dog. Error: %v", err)
	}
	return &dog, nil
}
