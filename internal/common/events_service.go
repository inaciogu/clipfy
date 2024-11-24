package common

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go-v2/service/sns/types"
	"github.com/oklog/ulid/v2"
	"os"
)

type EventsService struct {
	sns *sns.Client
}

type PublishMessageInput struct {
	Message        string
	MessageGroupId string
	Metadata       map[string]string
}

func NewEventsService(cfg aws.Config) *EventsService {
	return &EventsService{
		sns: sns.NewFromConfig(cfg),
	}
}

func (s *EventsService) Emit(input *PublishMessageInput) error {
	topicArn := os.Getenv("TOPIC_ARN")

	_, err := s.sns.Publish(context.TODO(), &sns.PublishInput{
		Message:                aws.String(input.Message),
		TopicArn:               aws.String(topicArn),
		MessageGroupId:         aws.String(input.MessageGroupId),
		MessageDeduplicationId: aws.String(ulid.Make().String()),
		MessageAttributes:      buildMetadata(input.Metadata),
	})

	if err != nil {
		return fmt.Errorf("failed to publish message: %v", err)
	}

	return nil
}

func buildMetadata(metadata map[string]string) map[string]types.MessageAttributeValue {
	if metadata == nil {
		return nil
	}
	var attributes map[string]types.MessageAttributeValue

	for key, value := range metadata {
		attributes[key] = types.MessageAttributeValue{
			DataType:    aws.String("String"),
			StringValue: aws.String(value),
		}
	}

	return attributes
}
