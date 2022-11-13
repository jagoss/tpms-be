package restclient

import (
	"fmt"
	"github.com/go-resty/resty/v2"
)

const (
	baseURL              = "localhost:8081"
	calculateVectorURL   = ""
	searchSimilarDogsURL = ""
	OK                   = 200
)

type CVModelClient struct {
	rc *resty.Client
}

func NewCVModelRestClient(client *resty.Client) *CVModelClient {
	return &CVModelClient{rc: client}
}

func (c *CVModelClient) CalculateVector(id int64, imgs []string) error {
	url := fmt.Sprintf("%s/%s", baseURL, calculateVectorURL)
	response, err := c.rc.R().
		SetHeader("Content-Type", "application/json").
		SetBody(CVRequest{DogID: id, Imgs: imgs}).
		Put(url)
	if err != nil {
		return err
	}
	if response.StatusCode() != OK {
		return fmt.Errorf("couldnt calculate vector for dog %d: %v", id, response.Error())
	}
	return nil
}

func (c *CVModelClient) SearchSimilarDog(dogID int64) ([]uint, error) {
	url := fmt.Sprintf("%s/%s/%d", baseURL, searchSimilarDogsURL, dogID)
	var resultList []uint
	response, err := c.rc.R().
		SetHeader("Content-Type", "application/json").
		SetResult(resultList).
		Get(url)
	if err != nil {
		return nil, err
	}
	if response.StatusCode() != OK {
		return nil, fmt.Errorf("couldnt get similar dogs for dog %d: %v", dogID, response.Error())
	}
	return resultList, nil
}

type CVRequest struct {
	DogID int64    `json:"dog_id"`
	Imgs  []string `json:"imgs"`
}
