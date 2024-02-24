package sys_updates

import (
	"context"
	"encoding/json"
	"fmt"
)

/* send notice change of status */

type NoticeStatusChange struct {
	Sref       string `json:"sref"`
	StatusCode int    `json:"status_code"`
}

func (nsc *NoticeStatusChange) Update(ctx context.Context) error {

	attr := make(map[string]string)
	attr["update"] = "notice_status_change"

	payload, err := json.Marshal(nsc)
	if err != nil {
		return fmt.Errorf("%w json NewIssuer", err)
	}

	return SendUpdate(ctx, payload, attr)
}
