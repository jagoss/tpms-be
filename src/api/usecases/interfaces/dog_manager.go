package interfaces

import "be-tpms/src/api/domain/model"

type DogManager interface {
	Get(string) (*model.Dog, error)
	Register(*model.Dog) (*model.Dog, error)
	Modify(*model.Dog) (*model.Dog, error)
	Delete(string) (bool, error)
}
