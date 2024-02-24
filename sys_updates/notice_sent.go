package sys_updates

import (
	"context"
	"encoding/json"
	"fmt"
)

/* Send notice has been sent */

type NoticeSent struct {
	Sref          string `json:"sref"`
	FleetId       int    `json:"fleet_id"`
	ScanForUpdate string `json:"scan_for_update"`
}

func (n *NoticeSent) Update(ctx context.Context) error {

	attr := make(map[string]string)
	attr["update"] = "parking_charge_notice"

	payload, err := json.Marshal(n)
	if err != nil {
		return fmt.Errorf("%w json NewIssuer", err)
	}

	return SendUpdate(ctx, payload, attr)
}
