package service

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"io"
)

type StorageService struct {
	client *s3.Client
	cfg    aws.Config
}

type UploadFileInput struct {
	FileName string    `json:"file_name"`
	File     io.Reader `json:"file"`
}

func NewS3Service(cfg aws.Config) *StorageService {
	return &StorageService{
		client: s3.NewFromConfig(cfg),
		cfg:    cfg,
	}
}

func (s *StorageService) UploadFile(input *UploadFileInput) error {
	_, err := s.client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String("clipfy-videos"),
		Key:    aws.String(input.FileName),
		Body:   input.File,
	})

	if err != nil {
		return fmt.Errorf("failed to upload file: %v", err)
	}

	return nil
}
