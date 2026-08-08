package searches

import (
	"cloud.google.com/go/firestore"
	"context"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/transfer360/sys360/registered_issuer"
	"github.com/transfer360/sys360/sys_updates"
	"google.golang.org/api/iterator"
)

//func searchZatparkSearch(ctx context.Context, sref string) {
//
//	client, err := firestore.NewClient(ctx, "transfer-360")
//
//	if err != nil {
//		log.Error("UpdateClientOnSearch:", err)
//		return err
//	}
//
//	defer client.Close()
//
//}

var ErrorSearchNotFound = errors.New("search not found")

func UpdateClientOnSearch(ctx context.Context, sref string, issuer registered_issuer.Issuer) (err error) {

	client, err := firestore.NewClient(ctx, "transfer-360")

	if err != nil {
		log.Error("UpdateClientOnSearch:", err)
		return err
	}

	defer client.Close()

	itr := client.Collection("searches").Where("sref", "==", sref).Where("result.is_hirer_vehicle", "==", true).Documents(ctx)

	documentID := ""

	for {
		doc, err := itr.Next()
		if err != nil {
			if errors.Is(err, iterator.Done) {
				break
			} else {
				log.Errorf("UpdateClientOnSearch iterating document. err=%v", err)
			}
		} else {
			documentID = doc.Ref.ID
			break
		}
	}

	if len(documentID) == 0 {
		return fmt.Errorf("%w [%s]", ErrorSearchNotFound, sref)
	}

	// Update the client fields individually: updating "client" as a whole would
	// delete every other key in the map, including software_id which downstream
	// services need to resolve the issuer's datetime strategy.
	updates := []firestore.Update{
		{Path: "client.clientid", Value: issuer.T360ID},
		{Path: "client.issuer_id", Value: issuer.T360ID},
	}

	if issuer.SoftwareProvider > 0 {
		updates = append(updates, firestore.Update{Path: "client.software_id", Value: issuer.SoftwareProvider})
	}

	_, err = client.Collection("searches").Doc(documentID).Update(ctx, updates)

	if err != nil {
		log.Error("UpdateClientOnSearch:", err)
		return err
	}

	// ------------------------------------------------------------------
	// Update sys_update
	// ------------------------------------------------------------------

	newi := sys_updates.ChangeIssuer{
		Sref:     sref,
		ClientID: issuer.T360ID,
	}

	err = newi.Update(ctx, fmt.Sprintf("sys360:UpdateClientOnSearch:%s", sref))
	if err != nil {
		log.Warnln("Pushing Packet for sys_updates:", err)
	}

	// ------------------------------------------------------------------

	return nil

}
