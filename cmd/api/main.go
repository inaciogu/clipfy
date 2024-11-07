package main

import (
	"clipfy/internal/api/command"
	"clipfy/internal/api/handler"
	"clipfy/internal/api/service"
	"context"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"
)

var ginLambda *ginadapter.GinLambda

func init() {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatalf("failed to load configuration, %v", err)
	}
	storageService := service.NewS3Service(cfg)
	uploadCommand := command.NewUploadFileCommand(storageService)
	uploadHandler := handler.NewUploadHandler(uploadCommand)

	log.Printf("Gin cold start")
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong1",
		})
	})
	r.POST("/upload", func(c *gin.Context) {
		uploadHandler.Handle(c)
	})

	if os.Getenv("ENV") == "PRODUCTION" {
		ginLambda = ginadapter.New(r)
		return
	}
	r.Run(":8080")
}

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// If no name is provided in the HTTP request body, throw an error
	return ginLambda.ProxyWithContext(ctx, req)
}

func main() {
	if os.Getenv("ENV") == "PRODUCTION" {
		lambda.Start(Handler)
	}

}
