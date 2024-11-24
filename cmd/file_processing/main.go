package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"os"
	"os/exec"
	"strings"
)

type EventBody struct {
	ID               string
	FileURL          string
	WithSubtitles    bool
	SegmentsDuration int64
	status           string
	UserID           string
}

func Handler(ctx context.Context, event events.SQSEvent) error {
	for _, record := range event.Records {
		var body EventBody
		err := json.Unmarshal([]byte(record.Body), &body)
		if err != nil {
			fmt.Println("Error unmarshalling body", err)
			return err
		}

		outputDir, err := os.MkdirTemp("", "segments-*")
		if err != nil {
			fmt.Println("Error creating temp dir", err)
			return err
		}

		fmt.Println("Output dir", outputDir)
		fmt.Println("Processing job", body.ID)
		fmt.Println("File URL", body.FileURL)

		// get file extension from url
		// split by "/"
		splitted := strings.Split(body.FileURL, "/")
		// get last element
		filename := splitted[len(splitted)-1]
		ext := strings.Split(filename, ".")[1]

		fmt.Println("File extension", ext)

		// Process the job
		cmd := exec.Command("/opt/bin/ffmpeg", "-i", fmt.Sprintf("https://%s", body.FileURL), "-c", "copy", "-map", "0", "-f", "segment", "-segment_time", fmt.Sprintf("%d", body.SegmentsDuration), "-reset_timestamps", "1", fmt.Sprintf("%s/segment_%%03d.%s", outputDir, ext))

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
		for _, file := range files {
			fmt.Println(file.Name())
		}

	}
	return nil
}

func main() {
	lambda.Start(Handler)
}
