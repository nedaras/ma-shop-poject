package apis

import (
	"errors"
	"nedas/shop/pkg/utils"
	"time"
)

var (
	ErrNotFound    = errors.New("address not found")
	ErrRateLimited = errors.New("rate limited")
	vlidator       AddressValidator
)

func ValidateAddress(address Address) (Address, error) {
	utils.Assert(vlidator != nil, "address validator is not set")
	return vlidator.ValidateAddress(address)
}

func GetTimeTillNextRequest() time.Duration {
	return vlidator.GetTimeTillNextRequest()
}

func SetAddressValidator(v AddressValidator) {
	utils.Assert(vlidator == nil, "address validator is already set")
	vlidator = v
}

type AddressValidator interface {

	// Any returned error should be of type [*AddressValidationError].
	ValidateAddress(adress Address) (Address, error)

	GetTimeTillNextRequest() time.Duration
}

type Address struct {
	Country string
	Street  string
	Region  string
	City    string
	Zipcode string
}

func (a Address) String() string {
	return a.Street + ", " + a.City + ", " + a.Region + ", " + a.Country + ", " + a.Zipcode
}

type AddressValidationError struct {
	Address Address
	Err     error
}

func (e *AddressValidationError) Error() string {
	return "'" + e.Address.String() + "': " + e.Err.Error()
}

func (e *AddressValidationError) Unwrap() error {
	return e.Err
}
