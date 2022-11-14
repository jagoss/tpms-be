package restclient

import (
	"be-tpms/src/api/domain/model"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

const (
	baseURL              = "http://161.35.228.212:8501"
	calculateEmbedding   = "/v1/models/model:predict"
	searchSimilarDogsURL = "/get_neighbors"
	OK                   = 200
)

type CVModelClient struct {
	rc *http.Client
}

func NewCVModelRestClient(client *http.Client) *CVModelClient {
	return &CVModelClient{rc: client}
}

func (c *CVModelClient) CalculateDogEmbedding(tensor model.Tensor) ([]float64, error) {
	body := map[string]interface{}{
		"instances": tensor.Values,
	}
	res, err := Post(fmt.Sprintf("%s%s", baseURL, calculateEmbedding), body)
	if err != nil {
		msg := fmt.Sprintf("[cvmodelrestclient.CalculateEmbedding] %s", err.Error())
		log.Printf(msg)
		return nil, fmt.Errorf(msg)
	}
	return res, nil
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
	ID    int64  `json:"id"`
	Image string `json:"image"`
}

func (c *CVModelClient) put(url string, reqBody *CVRequest) error {
	reqBodyJson, _ := json.Marshal(reqBody)
	request, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(reqBodyJson))
	request.Header.Set("Content-Type", "application/json")
	if err != nil {
		return err
	}
	res, err := http.DefaultClient.Do(request)
	if res.StatusCode != OK {
		return fmt.Errorf("status code %d: %v", res.StatusCode, res.Status)
	}
	return nil
}

func Post(url string, body interface{}) ([]float64, error) {
	log.Printf("request body: %v", body)
	reqBodyJson, _ := json.Marshal(body)
	response, err := http.Post(url, "application/json", bytes.NewBuffer(reqBodyJson))
	if err != nil {
		log.Printf(err.Error())
		return nil, err
	}
	bytesRes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	_ = response.Body.Close()

	if response.StatusCode != OK {
		var respBody map[string][]interface{}
		_ = json.Unmarshal(bytesRes, &respBody)
		return nil, fmt.Errorf("status code %s: %v", response.Status, respBody)
	}

	var respBody map[string][]float64
	_ = json.Unmarshal(bytesRes, &respBody)
	log.Printf("response body: %v", respBody)

	return respBody["predictions"], nil
}
