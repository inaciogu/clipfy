package main

import (
	"clipfy/internal/api/command"
	"clipfy/internal/api/handler"
	"clipfy/internal/api/service"
	"context"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/gin-gonic/gin"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/awslabs/aws-lambda-go-api-proxy/gin"
)

var ginLambda *ginadapter.GinLambda
var awsCfg aws.Config

func init() {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatalf("failed to load configuration, %v", err)
	}

	awsCfg = cfg
}

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// If no name is provided in the HTTP request body, throw an error
	return ginLambda.ProxyWithContext(ctx, req)
}

func main() {
	r := gin.Default()
	storageService := service.NewS3Service(awsCfg)
	uploadCommand := command.NewUploadFileCommand(storageService)
	uploadHandler := handler.NewUploadHandler(uploadCommand)

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
		lambda.Start(Handler)
		return
	}

	r.Run(":8080")
}
