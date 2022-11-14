package interfaces

import "be-tpms/src/api/io/restclient"

type CVModelRestClient interface {
	//CalculateEmbedding(id int64, imgs []string) error
	CalculateEmbedding() (*restclient.Tensor, error)
	SearchSimilarDog(dogID int64) ([]uint, error)
}
