package interfaces

import (
	"be-tpms/src/api/domain/model"
)

type DogPersiter interface {
	InsertDog(dog *model.User) (*model.User, error)
	GetDog(dogID uint) (*model.Dog, error)
	UpdateDog(dog *model.Dog) (*model.Dog, error)
	DeleteDog(userID string) error
}
