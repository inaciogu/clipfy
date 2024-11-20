package model

// job status enum
type JobStatus string

const (
	JobStatusPending    JobStatus = "PENDING"
	JobStatusProcessing JobStatus = "PROCESSING"
	JobStatusCompleted  JobStatus = "COMPLETED"
	JobStatusFailed     JobStatus = "FAILED"
)

type EditionJob struct {
	ID               string
	FileURL          string
	Status           JobStatus
	SegmentsDuration int64
	WithSubtitles    bool
	UserId           string
}
