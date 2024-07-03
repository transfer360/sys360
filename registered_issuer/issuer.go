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
	"google.golang.org/api/iterator"
	"strings"
	"time"
)

var ErrSavingNewIssuer = errors.New("unexpected issue saving the issuer")
var ErrConnectingToDatabase = errors.New("unexpected issue connecting to the database")
var ErrReadingFromDatabase = errors.New("unexpected issue reading from the database")
var ErrIssuerNotFound = errors.New("the registered issuer is not found")

const COLLECTION_REGISTERED_ISSUERS = "registered_issuers"

type Issuer struct {
	Issuer           string    `json:"issuer" firestore:"issuer"`
	PrivateParking   bool      `json:"private_parking" firestore:"private_parking"`
	Registered       time.Time `json:"-" firestore:"registered"`
	SoftwareProvider int       `json:"-" firestore:"software_provider"`
	Status           int       `json:"-" firestore:"status"`
	T360ID           string    `json:"t360_id" firestore:"t360_id"`
}

func (i *Issuer) Get(ctx context.Context, issuerID string) error {

	client, err := firestore.NewClient(ctx, "transfer-360")

	if err != nil {
		log.Error("Issuer:Get:", err)
		return ErrConnectingToDatabase
	}

	itr := client.Collection(COLLECTION_REGISTERED_ISSUERS).Where("t360_id", "==", issuerID).Documents(ctx)

	docRef := ""
	for {
		doc, err := itr.Next()
		if errors.Is(err, iterator.Done) {
			break
		}
		if doc.Exists() {
			err = doc.DataTo(&i)
			if err != nil {
				log.Errorf("Issuer:Get:[%s] %v", issuerID, err)
				return fmt.Errorf("%w with id: %s", ErrReadingFromDatabase, issuerID)
			}
			docRef = doc.Ref.ID
		}
	}

	if len(docRef) == 0 {
		return fmt.Errorf("%w with id: %s", ErrIssuerNotFound, issuerID)
	}

	return nil

}

func (i Issuer) Create(ctx context.Context, softwareProvider int) (Issuer, error) {

	ni := Issuer{}

	if len(i.Issuer) == 0 {
		return ni, fmt.Errorf("missing issuers name")
	}

	ni.Issuer = i.Issuer
	ni.PrivateParking = true
	ni.Registered = time.Now()
	ni.SoftwareProvider = softwareProvider
	ni.Status = 0

	var t360ID string
	counter := 0
	preFix := "T360"
	if softwareProvider == 2 {
		preFix = "ZP"
	}

	for {

		t360ID = ni.Issuer

		if len(ni.Issuer) >= 4 {
			t360ID = fmt.Sprintf("%s%s%d", preFix, strings.ReplaceAll(ni.Issuer[:4], " ", ""), counter)
		}

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

	_ = saveAPIKey(ctx, i)

	err = publish.Push(ctx, "transfer-360", "new_operator_registered", packData, nil)
	if err != nil {
		log.Error("Issuer:save:", err)
		return ErrSavingNewIssuer
	}

	return nil
}

func saveAPIKey(ctx context.Context, i Issuer) error {

	client, err := firestore.NewClient(ctx, "transfer-360")

	if err != nil {
		log.Error("saveAPIKey:", err)
		return ErrSavingNewIssuer
	}

	defer client.Close()

	data := struct {
		Active      bool   `firestore:"active"`
		ApiKey      string `firestore:"api_key"`
		ClientID    string `firestore:"client_id"`
		Description string `firestore:"description"`
		SoftwareID  int    `firestore:"software_id"`
	}{
		ApiKey:      uuid.New().String(),
		ClientID:    i.T360ID,
		Description: i.Issuer,
		SoftwareID:  i.SoftwareProvider,
	}

	_, err = client.Collection("api_keys").NewDoc().Set(ctx, data)
	if err != nil {
		log.Error("saveAPIKey:", err)
		return ErrSavingNewIssuer
	}

	return nil
}
