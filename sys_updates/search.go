package sys_updates

import (
	"context"
	"encoding/json"
	"fmt"
)

/* send result of search */

type SearchData struct {
	Sref                 string   `json:"sref"`
	ClientID             string   `json:"client_id"`
	IsHireVehicle        bool     `json:"is_hire_vehicle"`
	FleetID              int      `json:"fleet_id"`
	Reference            string   `json:"reference"`
	VRM                  string   `json:"vrm"`
	ScanDateTime         string   `json:"scan_date_time"`
	PercentSuccess       float64  `json:"percent_success"`
	IgnoreSearchPartners []string `json:"ignore_search_partners"`
	SearchErrors         []string `json:"search_errors"`
}

func (sd *SearchData) Update(ctx context.Context, source string) error {

	attr := make(map[string]string)
	attr["update"] = "search"
	attr["source"] = source

	payload, err := json.Marshal(sd)
	if err != nil {
		return fmt.Errorf("%w json NewIssuer", err)
	}

	return SendUpdate(ctx, payload, attr)
}
