package validate

import (
	"cloud.google.com/go/firestore"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"google.golang.org/api/iterator"
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

type BadAddress struct {
	CompanyName  string `firestore:"company_name"`
	AddressLine1 string `firestore:"addressline1"`
	AddressLine2 string `firestore:"addressline2"`
	AddressLine3 string `firestore:"addressline3"`
	AddressLine4 string `firestore:"addressline4"`
	Postcode     string `firestore:"postcode"`
}

var ErrAddressNotfound = errors.New("address not found")

func BadAddressLookUp(ctx context.Context, companyName string, PostCode string) (Address, error) {

	addr := Address{}
	baddr := BadAddress{}

	client, err := firestore.NewClient(ctx, "transfer-360")
	if err != nil {
		fmt.Errorf("Failed to create client: %v\n", err)
		return addr, err
	}
	// Ensure the context isn't canceled
	defer client.Close()

	var badAddressFound bool

	itr := client.Collection("bad_addresses").Where("company_name", "==", companyName).Where("postcode", "==", PostCode).Limit(1).Documents(ctx)
	for {
		doc, err := itr.Next()
		if err != nil {
			if errors.Is(err, iterator.Done) {
				break
			} else {
				log.Errorln("BadAddress: %v", err)
				return addr, err
			}
		} else {
			err = doc.DataTo(&baddr)
			if err != nil {
				log.Errorln("BadAddress: %v", err)
				return addr, err
			} else {
				badAddressFound = true
				break
			}
		}
	}

	if !badAddressFound {
		return addr, ErrAddressNotfound
	}

	addr.AddressLine1 = baddr.AddressLine1
	addr.AddressLine2 = baddr.AddressLine2
	addr.AddressLine3 = baddr.AddressLine3
	addr.City = baddr.AddressLine4
	addr.Postcode = baddr.Postcode

	if len(addr.City) == 0 {
		if len(addr.AddressLine3) > 0 {
			addr.City = addr.AddressLine3
			addr.AddressLine3 = ""
		} else if len(addr.AddressLine2) > 0 {
			addr.City = addr.AddressLine2
			addr.AddressLine2 = ""
		}
	}

	return addr, nil
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

	/*
		Redmond%20King%20County%20Durham%20SR8%202RR
	*/

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

				ukMapID := 0

				for i, ac := range mapInfo.Results {
					if ukMapID > 0 {
						break
					}
					for _, ac := range ac.AddressComponents {
						for _, art := range ac.Types {
							if art == "country" {
								if ac.ShortName == "GB" {
									ukMapID = i
									break
								}
							}
						}
					}
				}

				for _, ar := range mapInfo.Results[ukMapID].AddressComponents {
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

func AddressLineValidation(address string, postCode string) (Address, error) {

	addr := Address{}

	if len(os.Getenv("MAP_KEY")) == 0 {
		return addr, ErrMissingAddressKey
	}

	addressURL := strings.ReplaceAll(address, " ", "%20")

	mapsResponse, err := http.Get(fmt.Sprintf("https://maps.googleapis.com/maps/api/geocode/json?key=%s&address=%s", os.Getenv("MAP_KEY"), addressURL))

	if err != nil {
		log.Errorln("AddressLineValidation:", err)
		return addr, fmt.Errorf("%w %v", ErrUnexpected, err)
	} else {
		defer mapsResponse.Body.Close()

		mapInfo := data{}

		err = json.NewDecoder(mapsResponse.Body).Decode(&mapInfo)

		if err != nil {
			log.Errorln("AddressLineValidation:", err)
			return addr, fmt.Errorf("%w %v", ErrUnexpected, err)
		} else {

			if mapInfo.Status == "OK" {

				ukMapID := 0

				for i, ac := range mapInfo.Results {
					if ukMapID > 0 {
						break
					}
					for _, ac := range ac.AddressComponents {
						for _, art := range ac.Types {
							if art == "country" {
								if ac.ShortName == "GB" {
									ukMapID = i
									break
								}
							}
						}
					}
				}

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
