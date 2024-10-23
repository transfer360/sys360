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
	NoticeType    int    `json:"notice_type"`
	ScanForUpdate string `json:"scan_for_update"`
}

func (n *NoticeSent) Update(ctx context.Context, source string) error {

	updateType := ""

	if n.NoticeType == 0 || n.NoticeType == 1 {
		updateType = "parking_charge_notice"
	} else if n.NoticeType == 2 {
		updateType = "fuel_notice"
	} else {
		return fmt.Errorf("invalid notice type: %d", n.NoticeType)
	}

	attr := make(map[string]string)
	attr["update"] = updateType
	attr["source"] = source

	payload, err := json.Marshal(n)
	if err != nil {
		return fmt.Errorf("%w json NewIssuer", err)
	}

	return SendUpdate(ctx, payload, attr)
}
