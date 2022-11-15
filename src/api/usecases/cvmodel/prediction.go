package cvmodel

import (
	"be-tpms/src/api/domain/model"
	"be-tpms/src/api/io/db/persisters"
	"be-tpms/src/api/usecases/interfaces"
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"log"
	"math"
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
	var tensor [batch][width][height][depth]uint8
	for ix := 0; ix < width; ix++ {
		for iy := 0; iy < height; iy++ {
			img.ColorModel().Convert(img.At(ix, iy))
			r, g, b, _ := img.At(ix, iy).RGBA()
			tensor[0][ix][iy][0], tensor[0][ix][iy][1], tensor[0][ix][iy][2] = mapTo8bitValue(r), mapTo8bitValue(g), mapTo8bitValue(b)
		}
	}

	embedding, err := p.cvModel.CalculateDogEmbedding(model.Tensor{Values: tensor})
	if err != nil {
		msg := fmt.Sprintf("[prediction.CalculateEmbedding] error calculating embeding for dog %d: %s", dogID, err.Error())
		log.Printf(msg)
		return fmt.Errorf(msg)
	}

	err = p.dogPersister.UpdateEmbedding(dogID, fmt.Sprintf("%v", embedding))
	if err != nil {
		return err
	}

	return nil
}

func (p *Prediction) FindMatches(dogID uint) ([]model.Dog, error) {
	dog, err := p.dogPersister.GetDog(dogID)
	if err != nil {
		return nil, err
	}

	possibleMatchingDogs, err := p.dogPersister.GetPossibleMatchingDog(dog)
	if err != nil {
		return nil, err
	}

	log.Printf("Possible matching dogs from DB: %v", possibleMatchingDogs)

	if possibleMatchingDogs == nil || len(possibleMatchingDogs) == 0 {
		return make([]model.Dog, 0), nil
	}

	top5Dogs := top5Dogs(persisters.ToFloat64List(dog.Embedding), possibleMatchingDogs)

	log.Printf("Top 5 matching dogs: %v", possibleMatchingDogs)

	if len(top5Dogs) == 0 {
		return make([]model.Dog, 0), nil
	}

	dogs, _ := p.dogPersister.GetDogs(top5Dogs)

	return dogs, nil
}

func mapTo8bitValue(val uint32) uint8 {
	return uint8(val / (0x0100 + 1))
}

func top5Dogs(desireDogVector []float64, compareVectors []model.DogVector) []uint {
	topDogs := make([]DogSimilarity, 5)
	for _, vector := range compareVectors {
		addToTop(vector.ID, calculateDistance(desireDogVector, vector.Vector), topDogs)
	}

	return getIDList(topDogs)
}

func getIDList(dogs []DogSimilarity) []uint {
	resultList := make([]uint, len(dogs))
	for _, dog := range dogs {
		resultList = append(resultList, dog.DogID)
	}
	return resultList
}

func calculateDistance(vector1 []float64, vector2 []float64) float64 {
	distance := float64(0)
	for i := range vector1 {
		distance += math.Pow(vector1[i]+vector2[i], 2)
	}
	return math.Sqrt(distance)
}

func addToTop(id uint, distance float64, topDogs []DogSimilarity) {
	possibleMatch := DogSimilarity{id, distance}
	if len(topDogs) < 5 {
		topDogs = append(topDogs, possibleMatch)
	} else {
		var movingDog DogSimilarity
		movingDog = possibleMatch
		for _, topDog := range topDogs {
			if possibleMatch.Distance < topDog.Distance {
				tempVector := topDog
				topDog = movingDog
				movingDog = tempVector
			}
		}
	}
}

type DogSimilarity struct {
	DogID    uint
	Distance float64
}
