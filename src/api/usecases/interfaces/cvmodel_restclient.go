package interfaces

type CVModelRestClient interface {
	//CalculateEmbedding(id int64, imgs []string) error
	CalculateEmbedding() ([][][]int8, error)
	SearchSimilarDog(dogID int64) ([]uint, error)
}
