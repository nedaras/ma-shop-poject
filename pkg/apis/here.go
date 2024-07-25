package apis

import (
	"encoding/json"
	"errors"
	"fmt"
	"nedas/shop/pkg/utils"
	"net/http"
	"net/url"
)

// todo: we need to check the scoring and stuff...
type HereData struct {
	Items []struct {
		Address struct {
			CountryName string `json:"countryName"`
			State       string `json:"state"`
			City        string `json:"city"`
			Street      string `json:"street"`
			PostalCode  string `json:"postalCode"`
			HouseNumber string `json:"houseNumber"`
		} `json:"address"`
	} `json:"items"`
}

type Here struct{}

func (h *Here) ValidateAddress(address Address) (Address, error) {
	utils.Assert(address.Country != "", "country is empty")
	utils.Assert(address.Street != "", "address line is empty")
	utils.Assert(address.Region != "", "regionis is empty")
	utils.Assert(address.City != "", "city is empty")
	utils.Assert(address.Zipcode != "", "zipcode is empty")

	params := url.Values{}
	params.Add("q", address.String())
	params.Add("apiKey", utils.Getenv("HERE_API_KEY"))

	res, err := http.Get("https://geocode.search.hereapi.com/v1/geocode?" + params.Encode())
	if err != nil {
		return Address{}, &AddressValidationError{Address: address, Err: err}
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return Address{}, &AddressValidationError{Address: address, Err: fmt.Errorf("got unexpected response code '%d'", res.StatusCode)}
	}

	if res.Header.Get("Content-Type") != "application/json" {
		return Address{}, &AddressValidationError{Address: address, Err: errors.New("responded content is not in json form")}
	}

	data := &HereData{}
	decoder := json.NewDecoder(res.Body)

	if err := decoder.Decode(data); err != nil {
		return Address{}, &AddressValidationError{Address: address, Err: err}
	}

	if len(data.Items) == 0 {
		return Address{}, &AddressValidationError{Address: address, Err: ErrNotFound}
	}

	// todo: check if fields not null or sum
	return Address{
		Country: data.Items[0].Address.CountryName,
		Street:  data.Items[0].Address.Street + " " + data.Items[0].Address.HouseNumber,
		Region:  data.Items[0].Address.State,
		City:    data.Items[0].Address.City,
		Zipcode: data.Items[0].Address.PostalCode,
	}, nil
}
