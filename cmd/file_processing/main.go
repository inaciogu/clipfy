package main

import (
	"bytes"
	"clipfy/internal/api/command"
	"clipfy/internal/api/service"
	"clipfy/internal/common"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
)

type EventBody struct {
	ID               string
	FileURL          string
	WithSubtitles    bool
	SegmentsDuration int64
	status           string
	UserID           string
}

var storageService *common.StorageService
var segmentsService *service.SegmentsService

func init() {
	fmt.Println("Starting file_processing lambda")
	cfg, err := config.LoadDefaultConfig(context.Background())

	if err != nil {
		fmt.Println("Error loading AWS config", err)
		os.Exit(1)
	}

	storageService = common.NewS3Service(cfg)
	segmentsService = service.NewSegmentsService(cfg)
}

func Handler(ctx context.Context, event events.SQSEvent) error {
	for _, record := range event.Records {
		var body EventBody
		err := json.Unmarshal([]byte(record.Body), &body)
		if err != nil {
			fmt.Println("Error unmarshalling body", err)
			return err
		}

		outputDir := "/tmp"
		cdnURL := os.Getenv("CDN_URL")

		fmt.Println("Output dir", outputDir)
		fmt.Println("Processing job", body.ID)
		fmt.Println("File URL", body.FileURL)

		splitted := strings.Split(body.FileURL, "/")
		objectName := splitted[len(splitted)-1]
		fileParts := strings.Split(objectName, ".")
		ext := fileParts[len(fileParts)-1]
		fileName := fileParts[0]

		fmt.Println("File extension", ext)

		// Process the job
		cmd := exec.Command("/opt/bin/ffmpeg", "-i", fmt.Sprintf("https://%s", body.FileURL), "-c", "copy", "-map", "0", "-f", "segment", "-segment_time", fmt.Sprintf("%d", body.SegmentsDuration), "-reset_timestamps", "1", fmt.Sprintf("%s/%s_segment_%%03d.%s", outputDir, fileName, ext))

		var stdout, stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr
		err = cmd.Run()
		if err != nil {
			fmt.Println("Error processing job", stderr.String())
			return err
		}

		fmt.Println("Job processed successfully")
		files, err := os.ReadDir(outputDir)
		if err != nil {
			fmt.Println("Error reading output dir", err)
			return err
		}

		var segments []*command.CreateSegmentsInput
		for _, file := range files {
			fmt.Println("uploading file", file.Name())
			stream, err := os.Open(fmt.Sprintf("%s/%s", outputDir, file.Name()))
			if err != nil {
				fmt.Println("Error opening file", err)
				return err
			}

			err = storageService.UploadFile(ctx, os.Getenv("BUCKET_NAME"), fmt.Sprintf("%s/%s", body.UserID, file.Name()), stream)
			if err != nil {
				fmt.Println("Error uploading file", err)
				return err
			}

			stream.Close()

			segments = append(segments, &command.CreateSegmentsInput{
				ParentID:    body.ID,
				ParentName:  fileName,
				SegmentName: file.Name(),
				SegmentURL:  fmt.Sprintf("%s/%s/%s", cdnURL, body.UserID, file.Name()),
			})

			fmt.Println("File uploaded successfully")
		}

		_, err = command.NewCreateSegmentsCommand(segmentsService).Execute(segments)

		if err != nil {
			fmt.Println("Error creating segments", err)
			return err
		}

		err = storageService.DeleteFile(ctx, os.Getenv("BUCKET_NAME"), fmt.Sprintf("%s/%s", body.UserID, fileName))

		if err != nil {
			return err
		}
	}
	return nil
}

func main() {
	lambda.Start(Handler)
}
