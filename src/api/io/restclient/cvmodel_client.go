package restclient

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"strconv"
)

const (
	baseURL              = "https://dog-recognition-app-4l8w5.ondigitalocean.app/dog-recognition2"
	calculateEmbedding   = "/generate_embedding"
	searchSimilarDogsURL = "/get_neighbors"
	OK                   = 200
)

type CVModelClient struct {
	rc *resty.Client
}

func NewCVModelRestClient(client *resty.Client) *CVModelClient {
	client.BaseURL = baseURL
	return &CVModelClient{rc: client}
}

func (c *CVModelClient) CalculateEmbedding(id int64, imgs []string) error {
	response, err := c.rc.R().
		SetHeader("Content-Type", "application/json").
		SetBody(CVRequest{ID: id, Imgs: imgs[0]}).
		Put(calculateEmbedding)
	if err != nil {
		return err
	}
	if response.StatusCode() != OK {
		return fmt.Errorf("couldnt calculate vector for dog %d: %v", id, response.Error())
	}
	return nil
}

func (c *CVModelClient) SearchSimilarDog(dogID int64) ([]uint, error) {
	var resultList []uint
	response, err := c.rc.R().
		SetHeader("Content-Type", "application/json").
		SetQueryParam("dog_id", strconv.FormatInt(dogID, 10)).
		SetResult(resultList).
		Get(searchSimilarDogsURL)
	if err != nil {
		return nil, err
	}
	if response.StatusCode() != OK {
		return nil, fmt.Errorf("couldnt get similar dogs for dog %d: %v", dogID, response.Error())
	}
	return resultList, nil
}

type CVRequest struct {
	ID   int64  `json:"id"`
	Imgs string `json:"image"`
}
