package command

import "clipfy/internal/api/service"

type ListSegmentsOutput struct {
	ID          string `json:"id"`
	ParentID    string `json:"parent_id"`
	ParentName  string `json:"parent_name"`
	SegmentName string `json:"segment_name"`
	SegmentURL  string `json:"segment_url"`
}

type ListSegmentsCommand struct {
	service *service.SegmentsService
}

func NewListSegmentsCommand(service *service.SegmentsService) *ListSegmentsCommand {
	return &ListSegmentsCommand{
		service: service,
	}
}

func (c *ListSegmentsCommand) Execute(parentID string) ([]*ListSegmentsOutput, error) {
	segments, err := c.service.GetSegments(parentID)

	if err != nil {
		return nil, err
	}

	output := make([]*ListSegmentsOutput, 0)

	for _, segment := range segments {
		out := &ListSegmentsOutput{
			ID:          segment.ID,
			ParentID:    segment.ParentID,
			ParentName:  segment.ParentName,
			SegmentName: segment.SegmentName,
			SegmentURL:  segment.SegmentURL,
		}

		output = append(output, out)
	}

	return output, nil
}
