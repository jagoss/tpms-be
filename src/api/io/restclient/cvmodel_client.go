package restclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
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

//func (c *CVModelClient) CalculateEmbedding(id int64, imgs []string) error {
//	err := c.put(fmt.Sprintf("%s/%s", baseURL, calculateEmbedding), &CVRequest{ID: id, Image: imgs[0]})
//	if err != nil {
//		msg := fmt.Sprintf("[cvmodelrestclient.CalculateEmbedding] %s", err.Error())
//		log.Printf(msg)
//		return fmt.Errorf(msg)
//	}
//	return nil
//}

func (c *CVModelClient) CalculateEmbedding() (*Tensor, error) {
	const (
		width  = 224
		height = 224
	)
	s2 := make([]uint8, width*height*3)

	// ...read from hdf5...

	//to1D := func(x, y, z int) int {
	//    return (z * height * width) + (y * width) + x
	//}
	// display the fields
	fmt.Printf(":: size: length %v  capacity %v\n", len(s2), cap(s2))
	var vector [3][224][224]int
	//img := image.NewRGBA(image.Rect(0, 0, width, height))
	for iz := 0; iz < 3; iz++ {
		for ix := 0; ix < width; ix++ {
			for iy := 0; iy < height; iy++ {
				vector[iz][ix][iy] = rand.Intn(255)

				//ir := to1D(ix, iy, 0)
				//ig := to1D(ix, iy, 1)
				//ib := to1D(ix, iy, 2)
				//img.SetRGBA(ix, iy, color.RGBA{R: s2[ir], G: s2[ig], B: s2[ib], A: 255})
			}
		}
	}
	//body := map[string]interface{}{
	//	"instances": vector,
	//}
	//res, err := Post(fmt.Sprintf("%s%s", baseURL, calculateEmbedding), body)
	//if err != nil {
	//	msg := fmt.Sprintf("[cvmodelrestclient.CalculateEmbedding] %s", err.Error())
	//	log.Printf(msg)
	//	return nil, fmt.Errorf(msg)
	//}
	//log.Printf("%v", res)
	return &Tensor{vector}, nil
}

type Tensor struct {
	T [3][224][224]int
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

func Post(url string, body interface{}) ([]int8, error) {
	reqBodyJson, _ := json.Marshal(body)
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(reqBodyJson))
	request.Header.Set("Content-Type", "application/json")
	if err != nil {
		return nil, err
	}
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Printf(err.Error())
		return nil, err
	}
	if response.StatusCode != OK {
		return nil, fmt.Errorf("status code %s: %v", response.Status, response.Body)
	}

	//reqBodyJson, _ := json.Marshal(body)
	//response, err := http.Post(url, "application/json", bytes.NewBuffer(reqBodyJson))
	//if err != nil {
	//	return nil, err
	//}

	bytesRes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	_ = response.Body.Close()

	var respBody map[string][]int8
	_ = json.Unmarshal(bytesRes, &respBody)
	return respBody["predictions"], nil
}
