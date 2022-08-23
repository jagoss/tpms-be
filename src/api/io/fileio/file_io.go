package fileio

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/google/uuid"
	"image"
	"image/jpeg"
	"io/ioutil"
	"os"
	"strings"
)

const (
	dir       = "../../../../imgs"
	extension = ".jpeg"
)

func SaveImgs(imgs [][]byte) (string, error) {
	var f *os.File
	imgPaths := ""
	uuidNbr, _ := uuid.NewRandom()
	for i, imgByte := range imgs {
		img, _, err := image.Decode(bytes.NewReader(imgByte))
		path := fmt.Sprintf("%s/dog_%s_%d%s", dir, uuidNbr.String(), i, extension)
		f, err := os.Create(path)
		if err != nil {
			return "", fmt.Errorf("error creating file: %v", err)
		}

		opt := jpeg.Options{
			Quality: 90,
		}
		err = jpeg.Encode(f, img, &opt)
		if err != nil {
			return "", fmt.Errorf("error encoding image %s: %v", f.Name(), err)
		}

		imgPaths = fmt.Sprintf("%s;%s", imgPaths, f.Name())
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(f)

	return imgPaths, nil
}

func GetAllImgs() ([][]byte, error) {
	filesInfo, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	imgFileNames := ""
	for _, fileInfo := range filesInfo {
		imgFileNames = fmt.Sprintf("%s;%s/%s", imgFileNames, dir, fileInfo.Name())
	}
	imgFileNames = strings.Replace(imgFileNames, ";", "", 1)
	return GetImgs(imgFileNames)
}

func GetImgs(filePaths string) ([][]byte, error) {
	var f *os.File
	var imgsArray [][]byte
	for _, path := range strings.Split(filePaths, ";") {
		imgBytes, err := getSingleImg(f, path)
		if err != nil {
			return nil, err
		}

		imgsArray = append(imgsArray, imgBytes)
	}
	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	return imgsArray, nil
}

func getSingleImg(f *os.File, path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	fileInfo, _ := f.Stat()
	size := fileInfo.Size()
	imgBytes := make([]byte, size)

	buffer := bufio.NewReader(f)
	_, err = buffer.Read(imgBytes)
	if err != nil {
		return nil, err
	}

	return imgBytes, nil
}
