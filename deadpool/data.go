package deadpool

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"github.com/transfer360/sys360/ntk_to_lease"
	"time"
)

type Data struct {
	Sref    string    `json:"sref"`
	Created time.Time `json:"created"`
	Source  string    `json:"source"`
	Data    struct {
		OperatorName string `json:"operator_name"`
		Datetime     string `json:"datetime"`
		Amount       string `json:"amount"`
		Ref          string `json:"ref"`
		Vrm          string `json:"vrm"`
		Ntk          string `json:"ntk"`
		AppealUrl    string `json:"appeal_url"`
		PayUrl       string `json:"pay_url"`
		Site         string `json:"site"`
	} `json:"data"`
}

func (d Data) Convert(data []byte) ([]byte, error) {

	sr := ntk_to_lease.Data{}

	err := json.Unmarshal(data, &d)
	if err != nil {
		log.Errorf("Error while unmarshalling data: %s", err)
		return nil, err
	}

	sr.Operator = d.Data.OperatorName
	sr.Vrm = d.Data.Vrm
	sr.NTKUrl = d.Data.Ntk
	sr.Reference = d.Data.Ref
	sr.Site = d.Data.Site
	sr.Amount = d.Data.Amount
	sr.DateTime = d.Data.Datetime
	sr.AppealURL = d.Data.AppealUrl
	sr.PayURL = d.Data.PayUrl

	jsonData, err := json.Marshal(sr)
	if err != nil {
		log.Errorf("Error while marshalling data: %s", err)
		return nil, err
	}
	return jsonData, nil

}
