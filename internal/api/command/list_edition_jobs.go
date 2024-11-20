package command

import "clipfy/internal/api/service"

type ListEditionJobs struct {
	service *service.EditionJobService
}

type ListEditionJobsInput struct {
	UserID string
}

// output must be an array of EditionJob
