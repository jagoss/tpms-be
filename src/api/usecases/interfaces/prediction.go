package interfaces

import "be-tpms/src/api/domain/model"

type Prediction interface {
	CalculateEmbedding(dogID int64, imgsUrl string) error
	FindMatches(dogID uint, persister DogPersister) ([]model.Dog, error)
}
