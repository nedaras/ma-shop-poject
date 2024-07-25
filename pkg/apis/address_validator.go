package apis

import (
	"errors"
)

var (
	ErrNotFound = errors.New("address not found")
)

type AddressValidator interface {

	// Any returned error should be of type [*AddressValidationError].
	VaidateAddress(adress Address) (Address, error)
}

type Address struct {
	Country string
	Street  string
	Region  string
	City    string
	Zipcode string
}

// todo: bench this string function cuz i dont know how what '+' works does it like allocate a new string
// or go somehow idk puts it insidade builder its intresting
func (a Address) String() string {
	return a.Street + " " + a.City + ", " + a.Region + ", " + a.Country + ", " + a.Zipcode
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
