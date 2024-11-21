package service

import (
	"clipfy/internal/api/model"
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"strconv"
)

type EditionJobService struct {
	dynamodb *dynamodb.Client
}

func NewEditionJobService(cfg aws.Config) *EditionJobService {
	return &EditionJobService{
		dynamodb: dynamodb.NewFromConfig(cfg),
	}
}

func (e *EditionJobService) Create(editionJob *model.EditionJob) error {
	_, err := e.dynamodb.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String("clipfy"),
		Item: map[string]types.AttributeValue{
			"pk":               &types.AttributeValueMemberS{Value: fmt.Sprintf("EditionJob#%s", editionJob.UserId)},
			"sk":               &types.AttributeValueMemberS{Value: fmt.Sprintf("EditionJob#%s", editionJob.ID)},
			"id":               &types.AttributeValueMemberS{Value: editionJob.ID},
			"fileUrl":          &types.AttributeValueMemberS{Value: editionJob.FileURL},
			"status":           &types.AttributeValueMemberS{Value: string(editionJob.Status)},
			"userId":           &types.AttributeValueMemberS{Value: editionJob.UserId},
			"segmentsDuration": &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", editionJob.SegmentsDuration)},
			"withSubtitles":    &types.AttributeValueMemberBOOL{Value: editionJob.WithSubtitles},
		},
	})

	if err != nil {
		return err
	}

	return nil
}

func (e *EditionJobService) GetById(userId, id string) (*model.EditionJob, error) {
	output, err := e.dynamodb.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String("clipfy"),
		Key: map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{Value: fmt.Sprintf("EditionJob#%s", userId)},
			"sk": &types.AttributeValueMemberS{Value: fmt.Sprintf("EditionJob#%s", id)},
		},
	})

	if err != nil {
		return nil, err
	}

	segmentsDurationInt, err := strconv.ParseInt(output.Item["segmentsDuration"].(*types.AttributeValueMemberN).Value, 10, 32)
	if err != nil {
		return nil, err
	}
	editionJob := model.EditionJob{
		ID:               output.Item["id"].(*types.AttributeValueMemberS).Value,
		FileURL:          output.Item["fileUrl"].(*types.AttributeValueMemberS).Value,
		Status:           model.JobStatus(output.Item["status"].(*types.AttributeValueMemberS).Value),
		UserId:           output.Item["userId"].(*types.AttributeValueMemberS).Value,
		SegmentsDuration: segmentsDurationInt,
		WithSubtitles:    output.Item["withSubtitles"].(*types.AttributeValueMemberBOOL).Value,
	}

	return &editionJob, nil
}

func (e *EditionJobService) List(userId string) ([]model.EditionJob, error) {
	output, err := e.dynamodb.Query(context.TODO(), &dynamodb.QueryInput{
		TableName:              aws.String("clipfy"),
		KeyConditionExpression: aws.String("pk = :pk"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{Value: fmt.Sprintf("EditionJob#%s", userId)},
		},
	})

	if err != nil {
		return nil, err
	}

	var editionJobs []model.EditionJob
	for _, item := range output.Items {
		segmentsDurationInt, err := strconv.ParseInt(item["segmentsDuration"].(*types.AttributeValueMemberN).Value, 10, 32)
		if err != nil {
			return nil, err
		}
		editionJob := model.EditionJob{
			ID:               item["id"].(*types.AttributeValueMemberS).Value,
			FileURL:          item["fileUrl"].(*types.AttributeValueMemberS).Value,
			Status:           model.JobStatus(item["status"].(*types.AttributeValueMemberS).Value),
			UserId:           item["userId"].(*types.AttributeValueMemberS).Value,
			SegmentsDuration: segmentsDurationInt,
			WithSubtitles:    item["withSubtitles"].(*types.AttributeValueMemberBOOL).Value,
		}
		editionJobs = append(editionJobs, editionJob)
	}

	return editionJobs, nil
}

func (e *EditionJobService) UpdateStatus(userId, id string, status model.JobStatus) error {
	_, err := e.dynamodb.UpdateItem(context.TODO(), &dynamodb.UpdateItemInput{
		TableName: aws.String("clipfy"),
		Key: map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{Value: fmt.Sprintf("EditionJob#%s", userId)},
			"sk": &types.AttributeValueMemberS{Value: fmt.Sprintf("EditionJob#%s", id)},
		},
		UpdateExpression: aws.String("SET #status = :status"),
		ExpressionAttributeNames: map[string]string{
			"#status": "status",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":status": &types.AttributeValueMemberS{Value: string(status)},
		},
	})

	if err != nil {
		return err
	}

	return nil
}
