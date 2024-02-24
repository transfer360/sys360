package sys_updates

import (
	"context"
	"fmt"
	"github.com/transfer360/sys360/publish"
)

func SendUpdate(ctx context.Context, payload []byte, attr map[string]string) error {

	err := publish.Push(ctx, PROJECTID, "sys_updates", payload, attr)
	if err != nil {
		return fmt.Errorf("%w publish", err)
	}
	return nil

}
