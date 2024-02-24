package sys_updates

import (
	"context"
	"encoding/json"
	"fmt"
)

/* Send change of issuer */

type ChangeIssuer struct {
	Sref     string `json:"sref"`
	ClientID string `json:"client_id"`
}

func (ni *ChangeIssuer) Update(ctx context.Context) error {

	attr := make(map[string]string)
	attr["update"] = "change_issuer"

	payload, err := json.Marshal(ni)
	if err != nil {
		return fmt.Errorf("%w json NewIssuer", err)
	}

	return SendUpdate(ctx, payload, attr)
}
