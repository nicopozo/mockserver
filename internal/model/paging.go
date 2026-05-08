package model

type Paging struct {
	Total  int64  `json:"total"`
	Limit  int32  `json:"limit" minimum:"0" maximum:"1000" default:"30"`
	LastID string `json:"last_id,omitempty"`
}
