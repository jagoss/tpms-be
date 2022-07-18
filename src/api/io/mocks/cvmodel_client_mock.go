package mock

import (
	"be-tpms/src/api/domain/model"
)

type CVModelClient struct {
	dogResponses []model.DogResponse
	counter      int
	errResponse  []error
}

func (client *CVModelClient) SearchDog() (*model.DogResponse, error) {
	dog := client.dogResponses[client.counter]
	err := client.errResponse[client.counter]
	client.counter += 1
	return &dog, err
}
