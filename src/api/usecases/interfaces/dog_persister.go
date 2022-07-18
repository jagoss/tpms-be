package interfaces

import (
	"be-tpms/src/api/domain/model"
)

type DogPersister interface {
	InsertDog(dog *model.Dog) (*model.Dog, error)
	GetDog(dogID uint) (*model.Dog, error)
	UpdateDog(dog *model.Dog) (*model.Dog, error)
	DeleteDog(dogID uint) error
}
