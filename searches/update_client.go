package searches

import (
	"cloud.google.com/go/firestore"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/transfer360/sys360/publish"
	"github.com/transfer360/sys360/registered_issuer"
	"google.golang.org/api/iterator"
)

func UpdateClientOnSearch(ctx context.Context, sref string, issuer registered_issuer.Issuer) (err error) {

	client, err := firestore.NewClient(ctx, "transfer-360")

	if err != nil {
		log.Error("UpdateClientOnSearch:", err)
		return err
	}

	defer client.Close()

	clientInfo := struct {
		ClientID string `json:"client_id"`
		IssuerID string `json:"issuer_id"`
	}{
		ClientID: issuer.T360ID,
		IssuerID: issuer.Issuer,
	}

	itr := client.Collection("searches").Where("sref", "==", sref).Where("result.is_hirer_vehicle", "==", true).Documents(ctx)

	documentID := ""

	for {
		doc, err := itr.Next()
		if errors.Is(err, iterator.Done) {
			break
		}
		if doc.Exists() {
			documentID = doc.Ref.ID
			break
		}
	}

	if len(documentID) == 0 {
		return fmt.Errorf("search not found with Sref: %s", sref)
	}

	_, err = client.Collection("searches").Doc(documentID).Update(ctx, []firestore.Update{
		{
			Path:  "client",
			Value: clientInfo,
		},
	})

	if err != nil {
		log.Error("UpdateClientOnSearch:", err)
		return err
	}

	// ------------------------------------------------------------------
	// Update sys_update
	// ------------------------------------------------------------------

	attr := make(map[string]string)
	attr["update"] = "change_issuer"

	payload, err := json.Marshal(struct {
		Sref     string `json:"sref"`
		ClientID string `json:"client_id"`
	}{
		Sref:     sref,
		ClientID: clientInfo.ClientID,
	})

	if err != nil {
		log.Warnln("Creating Packet for sys_updates:", err)
	} else {

		err = publish.Push(ctx, "transfer-360", "sys_updates", payload, attr)

		if err != nil {
			log.Warnln("Pushing Packet for sys_updates:", err)
		}
	}
	// ------------------------------------------------------------------

	return nil

}
