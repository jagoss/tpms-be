package interfaces

import "be-tpms/src/api/domain/model"

type DogManager interface {
	Get(uint) (*model.Dog, error)
	Register(*model.Dog) (*model.Dog, error)
	Modify(*model.Dog) (*model.Dog, error)
	Delete(uint) (bool, error)
	ReportLostDog(*model.Dog) (*model.Dog, error)
}
