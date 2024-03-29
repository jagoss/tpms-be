package interfaces

import "be-tpms/src/api/domain/model"

type Prediction interface {
	CalculateEmbedding(dogID uint) error
	FindMatches(dogID uint) ([]model.Dog, error)
}
