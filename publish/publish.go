package publish

import (
	"cloud.google.com/go/pubsub"
	"context"
	"fmt"
)

func Push(ctx context.Context, project string, topicName string, payloadData []byte, attr map[string]string) error {

	client, err := pubsub.NewClient(ctx, project)
	if err != nil {
		return fmt.Errorf("there was an error creating pub/sub client: %w", err)
	}

	topic := client.Topic(topicName)

	msg := &pubsub.Message{
		Data:       payloadData,
		Attributes: attr,
	}

	if _, err := topic.Publish(ctx, msg).Get(ctx); err != nil {
		return fmt.Errorf("there was an error sending the payloadData: %w", err)
	}

	return nil
}
