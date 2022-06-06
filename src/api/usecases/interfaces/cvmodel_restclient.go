package interfaces

import "be-tpms/src/api/domain/model"

type CVModelRestClient interface {
	SearchDog() (*model.DogResponse, error)
}
