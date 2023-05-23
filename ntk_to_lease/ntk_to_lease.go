package ntk_to_lease

import (
	"cloud.google.com/go/pubsub"
	"context"
	"encoding/json"
	"fmt"
	joonix "github.com/joonix/log"
	log "github.com/sirupsen/logrus"
	"os"
	"time"
)

type Package struct {
	SearchReference string      `json:"sref"`
	Source          string      `json:"source"`
	Created         time.Time   `json:"created"`
	Data            interface{} `json:"data"`
}

func Send(ctx context.Context, ntk Package, attr map[string]string) error {

	if len(os.Getenv("DEVELOPMENT")) == 0 {
		log.SetFormatter(joonix.NewFormatter())
	}

	log.SetLevel(log.DebugLevel)

	client, err := pubsub.NewClient(ctx, "transfer-360")
	if err != nil {
		log.Errorf("SendNTKtoLease:newclient: %v", err)
		return fmt.Errorf("error creating pubsub client: %w", err)
	}

	data, err := json.Marshal(ntk)
	if err != nil {
		log.Errorf("SendNTKtoLease:data: %v", err)
		return fmt.Errorf("error creating data packet: %w", err)
	}

	topic := client.Topic("ntk_to_lease")

	msg := &pubsub.Message{
		Data:       data,
		Attributes: attr,
	}

	if _, err := topic.Publish(ctx, msg).Get(ctx); err != nil {
		log.Errorf("SendNTKtoLease:publish: %v", err)
		return fmt.Errorf("error sending data packet: %w", err)
	}

	return nil

}
