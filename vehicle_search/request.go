package vehicle_search

import "time"

type SearchRequest struct {
	Sref     string    `json:"sref"`
	Ref      string    `json:"ref"`
	VRM      string    `json:"vrm"`
	Date     string    `json:"date"`
	Source   string    `json:"source"`
	Received time.Time `json:"received"`
}
