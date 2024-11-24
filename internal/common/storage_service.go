package common

import (
	"bytes"
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"os"
)

type StorageService struct {
	client *s3.Client
	cfg    aws.Config
}

type GeneratePresignedUploadURLInput struct {
	FileName string `json:"file_name"`
}

type GeneratePresignedUploadURLOutput struct {
	UploadURL string `json:"upload_url"`
	FileURL   string `json:"file_url"`
}

func NewS3Service(cfg aws.Config) *StorageService {
	return &StorageService{
		client: s3.NewFromConfig(cfg),
		cfg:    cfg,
	}
}

func (s *StorageService) UploadFile(ctx context.Context, bucket, key, contentType string, body []byte) error {
	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(key),
		ContentType: aws.String(contentType),
		Body:        bytes.NewReader(body),
	})

	if err != nil {
		return fmt.Errorf("failed to upload file: %v", err)
	}

	return nil
}

func (s *StorageService) GeneratePresignedUploadURL(input *GeneratePresignedUploadURLInput) (*GeneratePresignedUploadURLOutput, error) {
	presigner := s3.NewPresignClient(s.client)

	presigned, err := presigner.PresignPutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String("clipfy-videos"),
		Key:    aws.String(input.FileName),
	})

	if err != nil {
		return nil, fmt.Errorf("failed generation upload URL: %v", err)
	}

	return &GeneratePresignedUploadURLOutput{
		UploadURL: presigned.URL,
		FileURL:   buildURL(input.FileName),
	}, nil
}

func buildURL(key string) string {
	cdnURL := os.Getenv("CDN_URL")
	return fmt.Sprintf("%s/%s", cdnURL, key)
}
