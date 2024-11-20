package service

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go-v2/service/sns/types"
	"os"
)

type EventsService struct {
	sns *sns.Client
}

type PublishMessageInput struct {
	Message  string
	Metadata map[string]string
}

func NewEventsService(sns *sns.Client) *EventsService {
	return &EventsService{
		sns: sns,
	}
}

func (s *EventsService) Emit(input *PublishMessageInput) error {
	topicArn := os.Getenv("TOPIC_ARN")

	_, err := s.sns.Publish(context.TODO(), &sns.PublishInput{
		Message:           &input.Message,
		TopicArn:          aws.String(topicArn),
		MessageAttributes: buildMetadata(input.Metadata),
	})

	if err != nil {
		return fmt.Errorf("failed to publish message: %v", err)
	}

	return nil
}

func buildMetadata(metadata map[string]string) map[string]types.MessageAttributeValue {
	var attributes map[string]types.MessageAttributeValue

	for key, value := range metadata {
		attributes[key] = types.MessageAttributeValue{
			DataType:    aws.String("String"),
			StringValue: aws.String(value),
		}
	}

	return attributes
}
