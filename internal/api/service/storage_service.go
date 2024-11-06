package service

import (
	"bytes"
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"net/http"
	"os"
	"time"
)

type StorageService struct {
	client *s3.Client
	cfg    aws.Config
}

type UploadFileInput struct {
	FileName      string `json:"file_name"`
	File          []byte `json:"file"`
	ContentType   string `json:"content_type"`
	ContentLength int64  `json:"content_length"`
}

type UploadFileOutput struct {
	FileURL string `json:"file_url"`
}

func NewS3Service(cfg aws.Config) *StorageService {
	return &StorageService{
		client: s3.NewFromConfig(cfg),
		cfg:    cfg,
	}
}

func (s *StorageService) UploadFile(input *UploadFileInput) (*UploadFileOutput, error) {
	presignClient := s3.NewPresignClient(s.client)
	bucketName := os.Getenv("BUCKET_NAME")

	presignResponse, err := presignClient.PresignPostObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(input.FileName),
	}, func(options *s3.PresignPostOptions) {
		options.Expires = 5 * time.Minute
	})
	if err != nil {
		return nil, err
	}

	fileToUpload := bytes.NewReader(input.File)

	req, err := http.NewRequest(http.MethodPost, presignResponse.URL, fileToUpload)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", input.ContentType)
	req.ContentLength = input.ContentLength

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, err
	}

	return &UploadFileOutput{
		FileURL: buildObjectURL(input.FileName),
	}, nil
}

func buildObjectURL(key string) string {
	cdnUrl := os.Getenv("CDN_URL")
	return fmt.Sprintf("%s/%s", cdnUrl, key)
}
