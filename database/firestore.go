package database

import (
	"cloud.google.com/go/firestore"
	"context"
	log "github.com/sirupsen/logrus"
)

func FirestoreConnect(ctx context.Context) (*firestore.Client, error) {

	client, err := firestore.NewClient(ctx, "transfer-360")

	if err != nil {
		log.Errorln(err)
		return nil, err
	}

	return client, nil

}
