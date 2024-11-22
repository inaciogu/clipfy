package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
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

	}
	return nil
}

func main() {
	lambda.Start(Handler)
}
