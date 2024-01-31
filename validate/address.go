package validate

import (
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"strings"
)

var ErrMissingAddressKey = errors.New("missing map address key")
var ErrInvalidMapResults = errors.New("invalid map results")
var ErrUnexpected = errors.New("unexpected error")

type addressParts struct {
	LongName  string   `json:"long_name"`
	ShortName string   `json:"short_name"`
	Types     []string `json:"types"`
}

type data struct {
	Results []struct {
		AddressComponents []addressParts `json:"address_components"`
		FormattedAddress  string         `json:"formatted_address"`
	} `json:"results"`
	Status string `json:"status"`
}

type Address struct {
	AddressLine1 string
	AddressLine2 string
	AddressLine3 string
	City         string
	Country      string
	Postcode     string
}

func AddressValidation(address []string, postCode string) (Address, error) {

	addr := Address{}

	if len(os.Getenv("MAP_KEY")) == 0 {
		return addr, ErrMissingAddressKey
	}

	addressURL := strings.Join(address, " ")
	addressURL = strings.ReplaceAll(addressURL, ",", "")
	addressURL = fmt.Sprintf("%s %s", addressURL, postCode)
	addressURL = strings.ReplaceAll(addressURL, " ", "%20")

	mapsResponse, err := http.Get(fmt.Sprintf("https://maps.googleapis.com/maps/api/geocode/json?key=%s&address=%s", os.Getenv("MAP_KEY"), addressURL))

	if err != nil {
		log.Errorln("AddressValidation:", err)
		return addr, fmt.Errorf("%w %v", ErrUnexpected, err)
	} else {
		defer mapsResponse.Body.Close()

		mapInfo := data{}

		err = json.NewDecoder(mapsResponse.Body).Decode(&mapInfo)

		if err != nil {
			log.Errorln("AddressValidation:", err)
			return addr, fmt.Errorf("%w %v", ErrUnexpected, err)
		} else {

			if mapInfo.Status == "OK" {

				for _, ar := range mapInfo.Results[0].AddressComponents {
					for _, art := range ar.Types {

						switch art {
						case "street_number":
							addr.AddressLine1 = ar.LongName
							break
						case "route":
							addr.AddressLine1 = strings.TrimSpace(fmt.Sprintf("%s %s", addr.AddressLine1, ar.LongName))
							break
						case "locality":
							addr.AddressLine2 = ar.LongName
							break
						case "postal_town":
							addr.City = ar.LongName
							break
						case "administrative_area_level_2":
							if len(addr.City) == 0 {
								addr.City = ar.LongName
							}
							break
						case "country":
							addr.Country = ar.LongName
							break
						case "postal_code":
							addr.Postcode = ar.LongName
							break
						}

					}
				}

				// 28.11 - Override post code due to part post code being returned.
				if len(postCode) > 0 {
					addr.Postcode = postCode
				}

				return addr, nil

			} else {
				return addr, fmt.Errorf("%w %v", ErrInvalidMapResults, err)
			}
		}
	}
}
