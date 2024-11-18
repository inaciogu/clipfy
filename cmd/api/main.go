package main

import (
	"clipfy/internal/api/command"
	"clipfy/internal/api/handler"
	"clipfy/internal/api/middleware"
	"clipfy/internal/api/service"
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

var ginLambda *ginadapter.GinLambda

func init() {
	awsCfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatalf("failed to load configuration, %v", err)
	}

	router := gin.Default()
	storageService := service.NewS3Service(awsCfg)
	uploadCommand := command.NewUploadFileCommand(storageService)
	uploadHandler := handler.NewUploadHandler(uploadCommand)

	router.Use(middleware.AuthMiddleware())

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong1",
		})
	})
	router.POST("/upload", func(c *gin.Context) {
		uploadHandler.Handle(c)
	})

	ginLambda = ginadapter.New(router)
}

func Handler(ctx context.Context, req events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {
	// Converte o evento LambdaFunctionURLRequest para um request compatível com o Gin
	proxyReq := events.APIGatewayProxyRequest{
		HTTPMethod:            req.RequestContext.HTTP.Method,
		Path:                  req.RawPath,
		Headers:               req.Headers,
		Body:                  req.Body,
		IsBase64Encoded:       req.IsBase64Encoded,
		QueryStringParameters: req.QueryStringParameters,
	}

	// Chama o adaptador do Gin com o proxy request
	resp, err := ginLambda.ProxyWithContext(ctx, proxyReq)
	if err != nil {
		return events.LambdaFunctionURLResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       `{"message": "Internal Server Error"}`,
		}, nil
	}

	// Converte a resposta do proxy para o formato LambdaFunctionURLResponse
	return events.LambdaFunctionURLResponse{
		StatusCode: resp.StatusCode,
		Headers:    resp.Headers,
		Body:       resp.Body,
	}, nil
}

func main() {
	lambda.Start(Handler)
}
