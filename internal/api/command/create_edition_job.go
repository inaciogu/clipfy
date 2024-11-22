package command

import (
	"clipfy/internal/api/model"
	"clipfy/internal/api/service"
	"encoding/json"
	"fmt"
	"github.com/oklog/ulid/v2"
)

type CreateEditionJob struct {
	service *service.EditionJobService
	events  *service.EventsService
}

type CreateEditionJobInput struct {
	Subtitle         bool   `json:"subtitle"`
	SegmentsDuration int64  `json:"segments_duration"` // in seconds
	FileURL          string `json:"file_url"`
	UserID           string
}

type CreateEditionJobOutput struct {
	ID               string `json:"id"`
	SegmentsDuration int64  `json:"segments_duration"`
	FileURL          string `json:"file_url"`
	Status           string `json:"status"`
}

func NewCreateEditionJob(service *service.EditionJobService, events *service.EventsService) *CreateEditionJob {
	return &CreateEditionJob{
		service: service,
		events:  events,
	}
}

func (u *CreateEditionJob) Execute(input *CreateEditionJobInput) (*CreateEditionJobOutput, error) {
	id := ulid.Make().String()
	editionJob := &model.EditionJob{
		ID:               id,
		WithSubtitles:    input.Subtitle,
		SegmentsDuration: input.SegmentsDuration,
		FileURL:          input.FileURL,
		UserId:           input.UserID,
		Status:           model.JobStatusPending,
	}

	err := u.service.Create(editionJob)
	if err != nil {
		return nil, err
	}

	message, err := json.Marshal(editionJob)
	if err != nil {
		return nil, err
	}
	fmt.Println(string(message))

	err = u.events.Emit(&service.PublishMessageInput{
		Message:        string(message),
		MessageGroupId: input.UserID,
		Metadata:       nil,
	})
	if err != nil {
		return nil, err
	}

	return &CreateEditionJobOutput{
		ID:               editionJob.ID,
		SegmentsDuration: editionJob.SegmentsDuration,
		FileURL:          editionJob.FileURL,
		Status:           string(editionJob.Status),
	}, err
}
