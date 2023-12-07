package registered_issuer

import (
	"cloud.google.com/go/firestore"
	"context"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"google.golang.org/api/iterator"
)

const REGISTERED_OPERATOR_COLLECTION = "registered_issuers"

var ErrIssuerIsNotRegistered = errors.New("this issuer is not registered on transfer360")

func IsRegistered(ctx context.Context, operatorname string) (registeredUser Issuer, err error) {

	client, err := firestore.NewClient(ctx, "transfer-360")

	if err != nil {
		log.Error("IsRegistered:", err)
		return registeredUser, err
	}

	defer client.Close()

	itr := client.Collection(REGISTERED_OPERATOR_COLLECTION).Where("issuer", "==", operatorname).Documents(ctx)

	isserFound := false

	for {
		doc, err := itr.Next()

		if errors.Is(err, iterator.Done) {
			break
		}

		err = doc.DataTo(&registeredUser)
		if err != nil {
			log.Error("IsRegistered:", err)
			return registeredUser, err
		} else {
			isserFound = true
		}
	}

	if isserFound {
		return registeredUser, nil
	} else {
		return registeredUser, fmt.Errorf("%w - %s", ErrIssuerIsNotRegistered, operatorname)
	}
}
