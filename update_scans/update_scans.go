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
	Sref          string    `json:"sref,omitempty" bson:"sref"`
	ClientID      string    `json:"client_id,omitempty" bson:"client_id"`
	ScanForUpdate time.Time `json:"scan_for_update" bson:"scan_for_update"`
	Status        int       `json:"status,omitempty" bson:"status"`
	DocumentID    string    `json:"document_id,omitempty" bson:"document_id"`
	FleetID       int       `json:"fleet_id,omitempty" bson:"fleet_id"`
	NoticeType    int       `json:"notice_type,omitempty" bson:"notice_type"`
	Source        string    `json:"source,omitempty" bson:"source"`
}

func (d *Data) Send(ctx context.Context) error {

	if len(d.Sref) == 0 {
		return fmt.Errorf("missing sref")
	}
	if len(d.ClientID) == 0 {
		return fmt.Errorf("missing client_id")
	}
	if d.ScanForUpdate.IsZero() {
		return fmt.Errorf("missing scan_for_update")
	}
	if d.FleetID == 0 {
		return fmt.Errorf("missing fleet_id")
	}

	if d.NoticeType == 0 {
		d.NoticeType = 1
	}

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
