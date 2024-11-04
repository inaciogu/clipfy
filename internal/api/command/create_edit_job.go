package command

type CreateEditJobCommand struct {
}

type CreateEditJobCommandInput struct {
	Subtitle         bool   `json:"subtitle"`
	SegmentsDuration int    `json:"segments_duration"` // in seconds
	FileName         string `json:"file_name"`
}

func (u *CreateEditJobCommand) Execute(input *CreateEditJobCommandInput) {
	// save job in database
	// send event to SQS
}
