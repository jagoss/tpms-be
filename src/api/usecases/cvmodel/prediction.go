package cvmodel

import (
	"be-tpms/src/api/domain/model"
	"be-tpms/src/api/usecases/interfaces"
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"log"
)

type Prediction struct {
	dogPersister interfaces.DogPersister
	cvModel      interfaces.CVModelRestClient
	storage      interfaces.Storage
}

const (
	width  = 224
	height = 224
	depth  = 3
	batch  = 1
)

func NewPrediction(dogPersister interfaces.DogPersister, restClient interfaces.CVModelRestClient, storage interfaces.Storage) *Prediction {
	return &Prediction{
		dogPersister: dogPersister,
		cvModel:      restClient,
		storage:      storage,
	}
}

func (p *Prediction) CalculateEmbedding(dogID uint) error {
	dog, err := p.dogPersister.GetDog(dogID)
	if err != nil {
		return err
	}
	imgsUrl := dog.ImgUrl
	if imgsUrl == "" {
		return nil
	}
	imgs, err := p.storage.GetImgs(imgsUrl)
	if err != nil {
		return err
	}
	imgByte, _ := base64.StdEncoding.DecodeString(imgs[0])
	byteReader := bytes.NewReader(imgByte)
	img, _, err := image.Decode(byteReader)
	if err != nil {
		return err
	}
	var tensor [batch][width][height][depth]uint32
	for ix := 0; ix < width; ix++ {
		for iy := 0; iy < height; iy++ {
			tensor[0][ix][iy][0], tensor[0][ix][iy][1], tensor[0][ix][iy][2], _ = img.At(ix, iy).RGBA()
		}
	}

	embedding, err := p.cvModel.CalculateDogEmbedding(model.Tensor{Values: tensor})
	if err != nil {
		msg := fmt.Sprintf("[prediction.CalculateEmbedding] error calculating embeding for dog %d: %s", dogID, err.Error())
		log.Printf(msg)
		return fmt.Errorf(msg)
	}

	err = p.dogPersister.UpdateEmbedding(uint(dogID), fmt.Sprintf("%v", embedding))
	if err != nil {
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
