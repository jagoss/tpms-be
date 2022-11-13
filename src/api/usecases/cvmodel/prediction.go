package cvmodel

import (
	"be-tpms/src/api/domain/model"
	"be-tpms/src/api/usecases/interfaces"
)

type Prediction struct {
	cvModel interfaces.CVModelRestClient
	storage interfaces.Storage
}

func NewPrediction(restClient interfaces.CVModelRestClient, storage interfaces.Storage) *Prediction {
	return &Prediction{
		cvModel: restClient,
		storage: storage,
	}
}

func (p *Prediction) CalculateEmbedding(dogID int64, imgsUrl string) error {
	imgs, err := p.storage.GetImgs(imgsUrl)
	if err != nil {
		return err
	}

	if err = p.cvModel.CalculateEmbedding(dogID, imgs); err != nil {
		return err
	}

	return nil
}

func (p *Prediction) FindMatches(dogID uint, persister interfaces.DogPersister) ([]model.Dog, error) {
	similarDogs, err := p.cvModel.SearchSimilarDog(int64(dogID))
	if err != nil {
		return nil, err
	}
	if similarDogs == nil || len(similarDogs) == 0 {
		return make([]model.Dog, 0), nil
	}

	dogs, err := persister.GetDogs(similarDogs)
	if err != nil {
		return nil, err
	}

	return dogs, nil
}
