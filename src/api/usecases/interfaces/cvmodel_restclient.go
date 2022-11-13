package interfaces

type CVModelRestClient interface {
	CalculateEmbedding(id int64, imgs []string) error
	SearchSimilarDog(dogID int64) ([]uint, error)
}
