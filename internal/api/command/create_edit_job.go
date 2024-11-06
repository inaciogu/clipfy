package command

type CreateEditJobCommand struct {
}

type CreateEditJobCommandInput struct {
	Subtitle         bool   `json:"subtitle"`
	SegmentsDuration int    `json:"segments_duration"` // in seconds
	FileURL          string `json:"file_url"`
}

func (u *CreateEditJobCommand) Execute(input *CreateEditJobCommandInput) {
	// save job in database
	// send event to SQS
}
