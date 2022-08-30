package storage

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"image"
	"image/jpeg"
	"io/ioutil"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
)

const (
	Dir       = "../../../../imgs"
	extension = ".jpeg"
)

func SaveImgs(imgs [][]byte) (string, error) {
	var f *os.File
	imgPaths := ""
	uuidNbr, _ := uuid.NewRandom()
	for i, imgByte := range imgs {
		img, _, err := image.Decode(bytes.NewReader(imgByte))
		path := fmt.Sprintf("%s/dog_%s_%d%s", Dir, uuidNbr.String(), i, extension)
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

func DeleteImg(path string) error {
	return os.Remove(path)
}

func GetAllImgs() ([][]byte, error) {
	filesInfo, err := ioutil.ReadDir(Dir)
	if err != nil {
		return nil, err
	}
	imgFileNames := ""
	for _, fileInfo := range filesInfo {
		imgFileNames = fmt.Sprintf("%s;%s/%s", imgFileNames, Dir, fileInfo.Name())
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

func SaveTempImg(c *gin.Context, fileHeader *multipart.FileHeader) ([][]byte, string, error) {
	randName, _ := uuid.NewRandom()
	path := "./" + randName.String() + filepath.Ext(fileHeader.Filename)
	err := c.SaveUploadedFile(fileHeader, path)
	if err != nil {
		return nil, "", err
	}

	imgBuffer, err := GetImgs(path)
	if err != nil {
		return nil, "", err
	}

	return imgBuffer, path, nil
}
