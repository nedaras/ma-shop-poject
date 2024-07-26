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

// add logic for rate limiting like hold a request in a day we can have 1k requests so we need to calculate like how many requests can be handled
// from given time to 12h or we could like suffle with api keys, idk if that even legal but in sense we could get like 2k requests a day
func (h *Here) ValidateAddress(address Address) (Address, error) { // todo: like idk use multiple addresses
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

	if res.StatusCode != http.StatusOK { // todo: check for rate limiting
		if res.StatusCode == http.StatusTooManyRequests {
			return Address{}, &AddressValidationError{Address: address, Err: ErrRateLimited}
		}
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

	if len(data.Items) > 1 {
		// todo add logger or sum idk
		fmt.Println(&AddressValidationError{Address: address, Err: errors.New("got multiple addresses")})
	}

	item := data.Items[0]

	if item.Address.CountryName == "" {
		return Address{}, &AddressValidationError{Address: address, Err: ErrNotFound}
	}

	if item.Address.Street == "" {
		return Address{}, &AddressValidationError{Address: address, Err: ErrNotFound}
	}

	if item.Address.HouseNumber == "" {
		return Address{}, &AddressValidationError{Address: address, Err: ErrNotFound}
	}

	if item.Address.State == "" {
		return Address{}, &AddressValidationError{Address: address, Err: ErrNotFound}
	}

	if item.Address.City == "" {
		return Address{}, &AddressValidationError{Address: address, Err: ErrNotFound}
	}

	if item.Address.PostalCode == "" {
		return Address{}, &AddressValidationError{Address: address, Err: ErrNotFound}
	}

	return Address{
		Country: item.Address.CountryName,
		Street:  item.Address.Street + " " + item.Address.HouseNumber,
		Region:  item.Address.State,
		City:    item.Address.City,
		Zipcode: item.Address.PostalCode,
	}, nil
}
