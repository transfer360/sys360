package publish

import (
	"context"
	"fmt"
	"sync"

	"cloud.google.com/go/pubsub"
)

var (
	clients   = make(map[string]*pubsub.Client)
	clientsMu sync.RWMutex
)

func getClient(ctx context.Context, project string) (*pubsub.Client, error) {
	clientsMu.RLock()
	client, ok := clients[project]
	clientsMu.RUnlock()
	if ok {
		return client, nil
	}

	clientsMu.Lock()
	defer clientsMu.Unlock()

	// Double-check after acquiring write lock
	if client, ok = clients[project]; ok {
		return client, nil
	}

	var err error
	client, err = pubsub.NewClient(ctx, project)
	if err != nil {
		return nil, fmt.Errorf("error creating pub/sub client: %w", err)
	}
	clients[project] = client
	return client, nil
}

func Push(ctx context.Context, project string, topicName string, payloadData []byte, attr map[string]string) error {
	client, err := getClient(ctx, project)
	if err != nil {
		return err
	}

	topic := client.Topic(topicName)
	if _, err := topic.Publish(ctx, &pubsub.Message{
		Data:       payloadData,
		Attributes: attr,
	}).Get(ctx); err != nil {
		return fmt.Errorf("error sending the payloadData: %w", err)
	}
	return nil
}

func PushWithOrderingKey(ctx context.Context, project string, topicName string, payloadData []byte, attr map[string]string, OrderingKey string) error {
	client, err := getClient(ctx, project)
	if err != nil {
		return err
	}

	topic := client.Topic(topicName)
	topic.EnableMessageOrdering = true

	if _, err := topic.Publish(ctx, &pubsub.Message{
		Data:        payloadData,
		Attributes:  attr,
		OrderingKey: OrderingKey,
	}).Get(ctx); err != nil {
		return fmt.Errorf("error sending the payloadData: %w", err)
	}
	return nil
}
