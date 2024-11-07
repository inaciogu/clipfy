package service

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"io"
	"os"
)

type StorageService struct {
	client *s3.Client
	cfg    aws.Config
}

type UploadFileInput struct {
	FileName string    `json:"file_name"`
	File     io.Reader `json:"file"`
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
	_, err := s.client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(os.Getenv("BUCKET_NAME")),
		Key:    aws.String(input.FileName),
		Body:   input.File,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to upload file: %v", err)
	}

	return &UploadFileOutput{
		FileURL: buildObjectURL(input.FileName),
	}, nil
}

func buildObjectURL(key string) string {
	cdnUrl := os.Getenv("CDN_URL")
	return fmt.Sprintf("%s/%s", cdnUrl, key)
}
