package apis

import (
	"errors"
	"strings"
)

var (
	ErrNotFound = errors.New("address not found")
)

type AddressValidator interface {

	// Any returned error should be of type [*AddressValidationError].
	VaidateAddress(adress Address) (Address, error)
}

type Address struct {
	Country      string
	AddressLine1 string
	AddressLine2 string // optional
	Region       string
	City         string
	Zipcode      string
}

// todo: bench this string function cuz i dont know how what '+' works does it like allocate a new string
//
//	or go somehow idk puts it insidade builder its intresting
func (a Address) String() string {
	builder := strings.Builder{}
	builder.WriteString(a.AddressLine1 + " ")

	if a.AddressLine2 != "" {
		builder.WriteString(a.AddressLine2 + " ")
	}

	builder.WriteString(a.City + ", ")
	builder.WriteString(a.Region + ", ")
	builder.WriteString(a.Country + ", ")
	builder.WriteString(a.Zipcode)

	return builder.String()
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
