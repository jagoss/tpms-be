package storage

import (
	"bytes"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/uuid"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

type Storage struct {
	bucket     *s3.S3
	bucketName string
}

func NewBucket() *Storage {
	key := os.Getenv("ACCESS_KEY") // Access key pair. You can create access key pairs using the control panel or API.
	secret := os.Getenv("SPACE_SECRET")
	s3Config := &aws.Config{
		Credentials: credentials.NewStaticCredentials(key, secret, ""), // Specifies your credentials.
		Endpoint:    aws.String("https://nyc3.digitaloceanspaces.com"), // Find your endpoint in the control panel, under Settings. Prepend "https://".
		Region:      aws.String("nyc3"),                                // Must be "us-east-1" when creating new Spaces. Otherwise, use the region in your endpoint, such as "nyc3".
	}

	// Step 3: The new session validates your request and directs it to your Space's specified endpoint using the AWS SDK.
	newSession, err := session.NewSession(s3Config)
	if err != nil {
		return nil
	}
	return &Storage{s3.New(newSession), os.Getenv("BUCKET")}
}

func (s *Storage) SaveImgs(imgs [][]byte) (string, error) {
	imgPaths := ""
	uuidNbr, _ := uuid.NewRandom()
	for i, imgByte := range imgs {
		filePath := fmt.Sprintf("dog_%s_%d%s", uuidNbr.String(), i, extension)
		err := s.saveFile(filePath, imgByte)
		if err != nil {
			return "", err
		}
		imgPaths = fmt.Sprintf("%s;%s", imgPaths, filePath)
	}
	imgPaths = strings.Replace(imgPaths, ";", "", 1)
	return imgPaths, nil
}

func (s *Storage) GetAllImgsName() ([]string, error) {
	input := &s3.ListObjectsInput{
		Bucket: aws.String(s.bucketName),
	}

	objects, err := s.bucket.ListObjects(input)
	if err != nil {
		return nil, fmt.Errorf("[bucket.GetAllImgsName] error getting imgs: %v", err)
	}
	fileList := []string{""}
	for _, obj := range objects.Contents {
		fileList = append(fileList, aws.StringValue(obj.Key))
	}
	return fileList, nil
}

func (s *Storage) GetImgs(filePaths string) ([][]byte, error) {
	var buffArray [][]byte
	for _, path := range strings.Split(filePaths, ";") {
		buffImg, err := s.getFile(path)
		if err != nil {
			return nil, err
		}
		buffArray = append(buffArray, buffImg)
	}
	log.Printf("len buffArray: %d", len(buffArray))
	return buffArray, nil
}

func (s *Storage) saveFile(key string, imgBuffer []byte) error {
	reader := bytes.NewReader(imgBuffer)
	_, err := s.bucket.PutObject(
		&s3.PutObjectInput{
			Bucket:             aws.String(s.bucketName),
			Key:                aws.String(key),
			ACL:                aws.String("private"),
			Body:               aws.ReadSeekCloser(reader),
			ContentLength:      aws.Int64(int64(len(imgBuffer))),
			ContentType:        aws.String(http.DetectContentType(imgBuffer)),
			ContentDisposition: aws.String("attachment"),
		},
	)
	if err != nil {
		return fmt.Errorf("[bucket.saveFile] error saving into bucket %s: %v", s.bucket.Endpoint, err)
	}
	return nil
}

func (s *Storage) getFile(key string) ([]byte, error) {
	log.Printf("img path in storage: %s", key)
	result, err := s.bucket.GetObject(
		&s3.GetObjectInput{
			Bucket: aws.String(s.bucketName),
			Key:    aws.String(key),
		})
	if err != nil {
		return nil, fmt.Errorf("[bucket.GetFile] %s", err.Error())
	}
	buffImg, err := io.ReadAll(result.Body)
	if err != nil {
		return nil, fmt.Errorf("[bucket.GetFile] %s", err.Error())
	}
	return buffImg, nil
}