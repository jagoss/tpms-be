package restclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	baseURL              = "https://dog-recognition-app-4l8w5.ondigitalocean.app/dog-recognition2"
	calculateEmbedding   = "/generate_embedding"
	searchSimilarDogsURL = "/get_neighbors"
	OK                   = 200
)

type CVModelClient struct {
	rc *http.Client
}

func NewCVModelRestClient(client *http.Client) *CVModelClient {
	return &CVModelClient{rc: client}
}

func (c *CVModelClient) CalculateEmbedding(id int64, imgs []string) error {
	err := c.put(fmt.Sprintf("%s/%s", baseURL, calculateEmbedding), &CVRequest{ID: id, Imgs: imgs[0]})
	if err != nil {
		return fmt.Errorf("[cvmodelrestclient.CalculateEmbedding] %s", err.Error())
	}
	return nil
}

func (c *CVModelClient) SearchSimilarDog(dogID int64) ([]uint, error) {
	response, err := c.rc.Get(fmt.Sprintf("%s/%s", baseURL, searchSimilarDogsURL))

	if err != nil {
		return nil, fmt.Errorf("[cvmodelrestclient.SearchSimilarDog] %s", err.Error())
	}

	if response.StatusCode != OK {
		return nil, fmt.Errorf("[cvmodelrestclient.SearchSimilarDog] couldnt get similar dogs for dog %d: %v", dogID, err)
	}
	resultListByte, _ := io.ReadAll(response.Body)
	_ = response.Body.Close()

	var resultList []uint
	err = json.Unmarshal(resultListByte, &resultList)
	if err != nil {
		return nil, err
	}
	return resultList, nil
}

type CVRequest struct {
	ID   int64  `json:"id"`
	Imgs string `json:"image"`
}

func (c *CVModelClient) put(url string, body interface{}) error {
	byteBuffer := new(bytes.Buffer)
	_ = json.NewEncoder(byteBuffer).Encode(body)
	request, err := http.NewRequest(http.MethodPut, url, byteBuffer)
	if err != nil {
		return err
	}
	res, err := http.DefaultClient.Do(request)
	if res.StatusCode != OK {
		return fmt.Errorf("status code not 200. It is %d", res.StatusCode)
	}
	return nil
}
