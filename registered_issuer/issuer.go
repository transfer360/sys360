package registered_issuer

import (
	"cloud.google.com/go/firestore"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"github.com/transfer360/sys360/publish"
	"strings"
	"time"
)

var ErrSavingNewIssuer = errors.New("unexpected issue saving the issuer")

type Issuer struct {
	APIKey           string    `json:"-" firestore:"api_key"`
	Issuer           string    `json:"issuer" firestore:"issuer"`
	PrivateParking   bool      `json:"private_parking" firestore:"private_parking"`
	Registered       time.Time `json:"-" firestore:"registered"`
	SoftwareProvider int       `json:"-" firestore:"software_provider"`
	Status           int       `json:"-" firestore:"status"`
	T360ID           string    `json:"t360_id" firestore:"t360_id"`
}

func (i Issuer) Create(ctx context.Context, softwareProvider int) (Issuer, error) {

	ni := Issuer{}

	if len(i.Issuer) == 0 {
		return ni, fmt.Errorf("missing issuers name")
	}

	ni.Issuer = i.Issuer
	ni.APIKey = uuid.New().String()
	ni.PrivateParking = true
	ni.Registered = time.Now()
	ni.SoftwareProvider = softwareProvider
	ni.Status = 0

	var t360ID string
	counter := 0

	for {

		t360ID = fmt.Sprintf("T360%s%d", strings.ReplaceAll(ni.Issuer[:4], " ", ""), counter)

		exists, err := T360IdExists(ctx, t360ID)
		if err != nil {
			log.Error(err)
			return ni, fmt.Errorf("error trying to save issuer")
		}

		if !exists {
			break
		}

		counter++
	}

	ni.T360ID = t360ID

	return ni, ni.save(ctx)
}

// save - Save the issuer
func (i Issuer) save(ctx context.Context) error {
	client, err := firestore.NewClient(ctx, "transfer-360")

	if err != nil {
		log.Error("Issuer:save:", err)
		return ErrSavingNewIssuer
	}

	defer client.Close()

	_, err = client.Collection(REGISTERED_OPERATOR_COLLECTION).NewDoc().Set(ctx, i)

	if err != nil {
		log.Error("Issuer:save:", err)
		return ErrSavingNewIssuer
	}

	packData, err := json.Marshal(i)
	if err != nil {
		log.Error("Issuer:save:", err)
		return ErrSavingNewIssuer
	}

	err = publish.Push(ctx, "transfer-360", "new_operator_registered", packData, nil)
	if err != nil {
		log.Error("Issuer:save:", err)
		return ErrSavingNewIssuer
	}

	return nil
}
