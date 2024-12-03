package service

import (
	"clipfy/internal/api/model"
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/oklog/ulid/v2"
)

type SegmentsService struct {
	dynamo *dynamodb.Client
}

func NewSegmentsService(cfg aws.Config) *SegmentsService {
	return &SegmentsService{
		dynamo: dynamodb.NewFromConfig(cfg),
	}
}

func (s *SegmentsService) CreateSegments(segments []*model.Segment) error {
	var writeRequests []types.WriteRequest

	for _, segment := range segments {
		id := ulid.Make().String()
		segment.ID = id
		writeRequests = append(writeRequests, types.WriteRequest{
			PutRequest: &types.PutRequest{
				Item: map[string]types.AttributeValue{
					"pk":           &types.AttributeValueMemberS{Value: "Segments#" + segment.ParentID},
					"sk":           &types.AttributeValueMemberS{Value: "Segments#" + segment.ID},
					"id":           &types.AttributeValueMemberS{Value: segment.ID},
					"parent_id":    &types.AttributeValueMemberS{Value: segment.ParentID},
					"segment_url":  &types.AttributeValueMemberS{Value: segment.SegmentURL},
					"segment_name": &types.AttributeValueMemberS{Value: segment.SegmentName},
					"parent_name":  &types.AttributeValueMemberS{Value: segment.ParentName},
				},
			},
		})
	}

	_, err := s.dynamo.BatchWriteItem(context.TODO(), &dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]types.WriteRequest{
			"clipfy": writeRequests,
		},
	})

	return err
}

func (s *SegmentsService) GetSegments(parentID string) ([]*model.Segment, error) {
	output, err := s.dynamo.Query(context.TODO(), &dynamodb.QueryInput{
		TableName:              aws.String("clipfy"),
		KeyConditionExpression: aws.String("pk = :pk and begins_with(sk, :sk)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{Value: "Segments#" + parentID},
			":sk": &types.AttributeValueMemberS{Value: "Segments#"},
		},
	})

	if err != nil {
		return nil, err
	}

	var segments []*model.Segment

	for _, item := range output.Items {
		segment := &model.Segment{
			ID:          item["id"].(*types.AttributeValueMemberS).Value,
			ParentID:    item["parent_id"].(*types.AttributeValueMemberS).Value,
			ParentName:  item["parent_name"].(*types.AttributeValueMemberS).Value,
			SegmentName: item["segment_name"].(*types.AttributeValueMemberS).Value,
			SegmentURL:  item["segment_url"].(*types.AttributeValueMemberS).Value,
		}

		segments = append(segments, segment)
	}

	return segments, nil
}
