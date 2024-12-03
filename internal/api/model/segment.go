package model

type Segment struct {
	ID          string `json:"id"`
	ParentID    string `json:"parent_id"`
	ParentName  string `json:"parent_name"`
	SegmentName string `json:"segment_name"`
	SegmentURL  string `json:"segment_url"`
}
