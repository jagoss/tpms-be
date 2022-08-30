package storage

import (
	"strings"
	"testing"
)

func TestSaveImgs_Ok(t *testing.T) {
	correctPath := "./imgs/dog_"
	imgBytes := getTestDogImgByteArray()
	byteArray := [][]byte{imgBytes}
	actualPath, err := SaveImgs(byteArray)
	if err != nil {
		t.Fatalf("error should be nil but got: %v", err)
	}
	if !strings.Contains(actualPath, correctPath) {
		t.Fatalf("file path should be: %s\n but got: %s", correctPath, actualPath)
	}
}

func TestGetImgs_Ok(t *testing.T) {
	img, err := GetImgs("../../../../test_imgs/test_img_1.jpeg")

	if err != nil {
		t.Fatalf("error should be nil but got: %v", err)
	}

	if img == nil || img[0] == nil {
		t.Fatalf("byte array should not be nil")
	}
}

func TestGetAllImgs(t *testing.T) {
	_, err := GetAllImgs()
	if err != nil {
		t.Fatalf("error should be nil but got: %v", err)
	}
}

func getTestDogImgByteArray() []byte {
	img, _ := GetImgs("../../../../test_imgs/test_img_1.jpeg")
	return img[0]
}
