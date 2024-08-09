package apis

import (
	"encoding/json"
	"errors"
	"fmt"
	"nedas/shop/pkg/utils"
	"net/http"
	"net/url"
	"sync"
	"time"
)

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

type Here struct {
	mu          sync.Mutex
	lastRequest time.Time
	requests    uint32
	maxRequests uint32
}

func NewHere(maxRequests uint32) *Here {
	return &Here{
		maxRequests: maxRequests,
		requests:    0,
	}
}

func (h *Here) ValidateAddress(address Address) (Address, error) {
	utils.Assert(address.Country != "", "country is empty")
	utils.Assert(address.Street != "", "address line is empty")
	utils.Assert(address.Region != "", "regionis is empty")
	utils.Assert(address.City != "", "city is empty")
	utils.Assert(address.Zipcode != "", "zipcode is empty")

	params := url.Values{}
	params.Add("q", address.String())
	params.Add("apiKey", utils.Getenv("HERE_API_KEY"))

	if h.GetTimeTillNextRequest() > 0 {
		return Address{}, ErrRateLimited
	}

	res, err := http.Get("https://geocode.search.hereapi.com/v1/geocode?" + params.Encode())
	if err != nil {
		return Address{}, &AddressValidationError{Address: address, Err: err}
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		if res.StatusCode == http.StatusTooManyRequests {
			h.mu.Lock()
			utils.Logger().Warn(fmt.Sprintf("got rate limited by here when still have %d requests left", h.requests))
			h.mu.Unlock()

			return Address{}, &AddressValidationError{Address: address, Err: ErrRateLimited}
		}
		return Address{}, &AddressValidationError{Address: address, Err: fmt.Errorf("got unexpected response code '%d'", res.StatusCode)}
	}

	now := time.Now()
	h.mu.Lock()

	a := now.Truncate(time.Hour * 24)
	b := h.lastRequest.Truncate(time.Hour * 24)
	if a.Sub(b).Hours() >= 24 {
		h.requests = 0
	}

	h.requests++
	h.lastRequest = now
	h.mu.Unlock()

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
		utils.Logger().Warn("got multiple addresses", address, data.Items)
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

// test it
func (h *Here) GetTimeTillNextRequest() time.Duration {
	now := time.Now()

	beg := now.Truncate(time.Hour * 24)
	end := beg.Add(time.Hour * 24)

	h.mu.Lock()
	lastRequest := h.lastRequest
	requests := h.requests
	h.mu.Unlock()

	if (lastRequest == time.Time{}) {
		return 0
	}

	if lastRequest.Sub(beg) < 0 {
		return 0
	}

	timeLeft := end.Sub(now).Milliseconds()
	requestsLeft := h.maxRequests - requests
	rm := float64(requestsLeft) / float64(timeLeft)

	if float64(now.Sub(lastRequest).Milliseconds())*rm > 1.0 {
		return 0
	}

	nextRequest := lastRequest.Add(time.Millisecond * time.Duration(1.0/rm))
	return nextRequest.Sub(now)
}
