package command

import (
	"clipfy/internal/api/model"
	"clipfy/internal/api/service"
)

type CreateSegmentsInput struct {
	ParentID    string `json:"parent_id"`
	ParentName  string `json:"parent_name"`
	SegmentName string `json:"segment_name"`
	SegmentURL  string `json:"segment_url"`
}

type CreateSegmentsOutput struct {
	ParentName  string `json:"parent_name"`
	SegmentName string `json:"segment_name"`
	SegmentURL  string `json:"segment_url"`
}

type CreateSegmentsCommand struct {
	service *service.SegmentsService
}

func NewCreateSegmentsCommand(service *service.SegmentsService) *CreateSegmentsCommand {
	return &CreateSegmentsCommand{
		service: service,
	}
}

func (c *CreateSegmentsCommand) Execute(input []*CreateSegmentsInput) ([]*CreateSegmentsOutput, error) {
	var segments []*model.Segment

	for _, in := range input {
		segment := &model.Segment{
			ParentID:    in.ParentID,
			ParentName:  in.ParentName,
			SegmentName: in.SegmentName,
			SegmentURL:  in.SegmentURL,
		}

		segments = append(segments, segment)
	}

	err := c.service.CreateSegments(segments)

	if err != nil {
		return nil, err
	}

	output := make([]*CreateSegmentsOutput, 0)

	for _, segment := range segments {
		out := &CreateSegmentsOutput{
			ParentName:  segment.ParentName,
			SegmentName: segment.SegmentName,
			SegmentURL:  segment.SegmentURL,
		}

		output = append(output, out)
	}

	return output, nil
}
