package interfaces

import (
	"be-tpms/src/api/domain/model"
)

type CVModelRestClient interface {
	CalculateDogEmbedding(tensor model.Tensor) ([]float64, error)
	SearchSimilarDog(dogID int64) ([]uint, error)
}
