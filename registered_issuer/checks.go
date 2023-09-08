package registered_issuer

import (
	"cloud.google.com/go/firestore"
	"context"
	"errors"
	log "github.com/sirupsen/logrus"
	"google.golang.org/api/iterator"
)

// T360IdExists ---------------------------------------------------------------------------
// Check the T360 ID is no on the system
func T360IdExists(ctx context.Context, t360Exists string) (idExists bool, err error) {
	client, err := firestore.NewClient(ctx, "transfer-360")

	if err != nil {
		log.Error("T360IdExists:", err)
		return false, err
	}

	defer client.Close()

	itr := client.Collection(REGISTERED_OPERATOR_COLLECTION).Where("t360_id", "==", t360Exists).Documents(ctx)

	for {
		doc, err := itr.Next()

		if errors.Is(err, iterator.Done) {
			break
		}

		idExists = doc.Exists()
	}

	return idExists, nil
}
