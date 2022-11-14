package interfaces

import (
	"be-tpms/src/api/domain/model"
)

type DogPersister interface {
	InsertDog(dog *model.Dog) (*model.Dog, error)
	GetDog(dogID uint) (*model.Dog, error)
	UpdateDog(dog *model.Dog) (*model.Dog, error)
	DeleteDog(dogID uint) error
	DogExisitsByNameAndOwner(string, string) (bool, error)
	GetMissingDogs() ([]model.Dog, error)
	GetDogsByUser(userID string) ([]model.Dog, error)
	GetDogs(dogs []uint) ([]model.Dog, error)
	SetLostDog(id uint, lat float64, lng float64) error
	UpdateEmbedding(dogID uint, embedding string) error
}
