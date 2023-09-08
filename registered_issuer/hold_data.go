package registered_issuer

import (
	"context"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"github.com/transfer360/sys360/ntk_to_lease"
	"github.com/transfer360/sys360/publish"
	"time"
)

type HoldDataPacket struct {
	Data   ntk_to_lease.Data `json:"data"`
	T360ID string            `json:"t360_id"`
	Posted time.Time         `json:"posted"`
}

func HoldData(ctx context.Context, data ntk_to_lease.Data, t360ID string) error {

	payload, err := json.Marshal(HoldDataPacket{
		Data:   data,
		T360ID: t360ID,
		Posted: time.Now(),
	})

	if err != nil {
		log.Error(err)
		return err
	}

	return publish.Push(ctx, "transfer-360", "ntk_pending_operator_registration", payload, nil)

}
