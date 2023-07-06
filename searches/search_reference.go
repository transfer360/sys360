package searches

import (
	"cloud.google.com/go/firestore"
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"google.golang.org/api/iterator"
)

func GetSearchReference(ctx context.Context, reference string, vrm string) (sref string, err error) {

	client, err := firestore.NewClient(ctx, "transfer-360")

	if err != nil {
		return sref, fmt.Errorf("error creating firestore client: %w", err)
	}

	defer client.Close()

	doci := client.Collection("searches").Where("result.your_reference", "==", reference).Where("result.vrm", "==", vrm).Documents(ctx)

	for {
		doc, err := doci.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Warnln(err)
		} else {

			sref = doc.Ref.ID
			break

		}

	}

	return sref, nil
}
