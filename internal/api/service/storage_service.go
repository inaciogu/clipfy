package service

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"io"
	"net/http"
)

type StorageService struct {
	client *s3.Client
	cfg    aws.Config
}

type UploadFileInput struct {
	FileName      string    `json:"file_name"`
	File          io.Reader `json:"file"`
	ContentLength int64     `json:"content_length"`
}

func NewS3Service(cfg aws.Config) *StorageService {
	return &StorageService{
		client: s3.NewFromConfig(cfg),
		cfg:    cfg,
	}
}

func (s *StorageService) UploadFile(input *UploadFileInput) error {
	presigner := s3.NewPresignClient(s.client)

	presigned, err := presigner.PresignPutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String("clipfy-videos"),
		Key:    aws.String(input.FileName),
	})

	if err != nil {
		return fmt.Errorf("failed generation upload URL: %v", err)
	}

	req, err := http.NewRequest(http.MethodPut, presigned.URL, input.File)
	req.ContentLength = input.ContentLength

	if err != nil {
		return fmt.Errorf("failed to create request for upload: %v", err)
	}

	client := &http.Client{}

	res, err := client.Do(req)

	if err != nil {
		return fmt.Errorf("error uploading file: %v", err)
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(res.Body)
		return fmt.Errorf("failed to upload file: %v - %s", res.Status, body)
	}

	return nil
}
