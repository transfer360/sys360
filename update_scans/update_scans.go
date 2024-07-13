package update_scans

import (
	"context"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/transfer360/sys360/publish"
	"time"
)

type Data struct {
	Sref          string    `json:"sref,omitempty"`
	ClientID      string    `json:"client_id,omitempty"`
	ScanForUpdate time.Time `json:"scan_for_update"`
	Status        int       `json:"status,omitempty"`
	DocumentID    string    `json:"document_id,omitempty"`
	FleetID       int       `json:"fleet_id,omitempty"`
	NoticeType    int       `json:"notice_type,omitempty"`
	Source        string    `json:"source,omitempty"`
}

func (d *Data) Send(ctx context.Context) error {

	payload, err := json.Marshal(d)
	if err != nil {
		log.Errorln("error marshaling data:", err)
		return fmt.Errorf("%w json marshal", err)
	}

	attr := make(map[string]string)

	err = publish.PushWithOrderingKey(ctx, "transfer-360", "update_scans", payload, attr, "update_scans_key")
	if err != nil {
		log.Errorln(err)
		return fmt.Errorf("%w publish", err)
	}
	return nil
}
