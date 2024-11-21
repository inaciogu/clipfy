package command

import (
	"clipfy/internal/api/service"
)

type ListEditionJobs struct {
	service *service.EditionJobService
}

type ListEditionJobsInput struct {
	UserID string
}

type ListEditionJobsOutput struct {
	ID               string `json:"id"`
	SegmentsDuration int64  `json:"segments_duration"`
	FileURL          string `json:"file_url"`
	Status           string `json:"status"`
}

func NewListEditionJobs(service *service.EditionJobService) *ListEditionJobs {
	return &ListEditionJobs{
		service: service,
	}
}

func (u *ListEditionJobs) Execute(input *ListEditionJobsInput) []*ListEditionJobsOutput {
	editionJobs, err := u.service.List(input.UserID)
	if err != nil {
		return nil
	}

	var editionJobsOutput []*ListEditionJobsOutput
	for _, editionJob := range editionJobs {
		editionJobsOutput = append(editionJobsOutput, &ListEditionJobsOutput{
			ID:               editionJob.ID,
			SegmentsDuration: editionJob.SegmentsDuration,
			FileURL:          editionJob.FileURL,
			Status:           string(editionJob.Status),
		})
	}

	return editionJobsOutput
}
