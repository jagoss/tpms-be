package restclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
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
	err := c.put(fmt.Sprintf("%s/%s", baseURL, calculateEmbedding), buildRequestBody(id, imgs[0]))
	if err != nil {
		msg := fmt.Sprintf("[cvmodelrestclient.CalculateEmbedding] %s", err.Error())
		log.Printf(msg)
		return fmt.Errorf(msg)
	}
	return nil
}

func (c *CVModelClient) SearchSimilarDog(dogID int64) ([]uint, error) {
	response, err := c.rc.Get(fmt.Sprintf("%s/%s", baseURL, searchSimilarDogsURL))

	if err != nil {
		msg := fmt.Sprintf("[cvmodelrestclient.SearchSimilarDog] %s", err.Error())
		log.Printf(msg)
		return nil, fmt.Errorf(msg)
	}

	if response.StatusCode != OK {
		msg := fmt.Sprintf("[cvmodelrestclient.SearchSimilarDog] couldnt get similar dogs for dog %d: %v", dogID, err)
		log.Printf(msg)
		return nil, fmt.Errorf(msg)
	}
	resultListByte, _ := io.ReadAll(response.Body)
	_ = response.Body.Close()

	var resultList []uint
	err = json.Unmarshal(resultListByte, &resultList)
	if err != nil {
		msg := fmt.Sprintf("[cvmodelrestclient.SearchSimilarDog] error unmarshalling body: %s", err.Error())
		log.Printf(msg)
		return nil, fmt.Errorf(msg)
	}
	return resultList, nil
}

func buildRequestBody(id int64, imgs string) map[string]interface{} {
	return map[string]interface{}{
		"id":    uint(id),
		"image": imgs,
	}
}

type CVRequest struct {
	ID   int64  `json:"id"`
	Imgs string `json:"image"`
}

func (c *CVModelClient) put(url string, body map[string]interface{}) error {
	log.Printf("request body: %v", body)
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
