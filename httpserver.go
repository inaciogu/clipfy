package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

const (
	BucketName = "clipfy-videos"
	Region     = "us-east-1"
)

func generatePresignedURL(objectKey string, expiration time.Duration) (string, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(Region))
	if err != nil {
		return "", fmt.Errorf("falha ao carregar configuração: %v", err)
	}

	client := s3.NewFromConfig(cfg)
	presigner := s3.NewPresignClient(client)

	req, err := presigner.PresignPutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(BucketName),
		Key:    aws.String(objectKey),
	}, s3.WithPresignExpires(expiration))
	if err != nil {
		return "", fmt.Errorf("falha ao criar URL assinada: %v", err)
	}

	return req.URL, nil
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	objectKey := "originals/" + time.Now().Format("20060102150405") + ".mp4"

	presignedURL, err := generatePresignedURL(objectKey, 15*time.Minute)
	if err != nil {
		http.Error(w, "Falha ao gerar URL assinada", http.StatusInternalServerError)
		log.Printf("Erro ao gerar URL assinada: %v", err)
		return
	}

	req, err := http.NewRequest(http.MethodPut, presignedURL, r.Body)

	if err != nil {
		http.Error(w, "Erro ao criar requisição para upload", http.StatusInternalServerError)
		log.Printf("Erro ao criar requisição para URL assinada: %v", err)
		return
	}

	req.ContentLength = r.ContentLength

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Erro ao enviar arquivo para S3", http.StatusInternalServerError)
		log.Printf("Erro ao enviar para URL assinada: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, "Falha no upload para S3", http.StatusInternalServerError)
		log.Printf("Erro de status no upload para S3: %v", resp.Status)
		return
	}

	fmt.Fprintf(w, "Arquivo enviado com sucesso para o S3!")
}

func main() {
	cmd := exec.Command("ffmpeg", "-i", "./lambda.mp4", "-c", "copy", "-map", "0", "-segment_time", "20", "-reset_timestamps", "1", "-f", "segment", "output%03d.mp4")
	err := cmd.Run()
	if err != nil {
		log.Fatalf("Erro ao dividir arquivo: %v", err)
	}
}
