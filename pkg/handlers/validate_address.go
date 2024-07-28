package handlers

import (
	"errors"
	"nedas/shop/pkg/apis"
	"nedas/shop/pkg/models"
	"net/http"
	"net/url"
	"strconv"

	"github.com/labstack/echo/v4"
)

type AddressData struct {
	CountryCode string
	Contact     string
	Phone       string
	Street      string
	Region      string
	City        string
	Zipcode     string
}

// pattern ^[A-Za-zÄÖÜäöüßĄČĘĖĮŠŲŪŽąčęėįšųūž ]+$
// if error return 400 err code and html to update addressData or sum
func HandlePutAddress(c echo.Context) error {
	session := getSession(c)
	storage := getStorage(c)

	if session == nil {
		return newHTTPError(http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 8)
	if err != nil {
		return newHTTPError(http.StatusBadRequest, "param 'id' is not valid uint8")
	}

	addressData, err := getAddressData(c)
	if err != nil {
		return err
	}

	// we need to sanatize and validate the shit out of this cuz what the user writes here will go to and database
	// would be crazy if an user writed in like 1k long names or idk phone number without numbers
	address, err := apis.ValidateAddress(apis.Address{
		Country: addressData.CountryCode,
		Street:  addressData.Street,
		Region:  addressData.Region,
		City:    addressData.City,
		Zipcode: addressData.Zipcode,
	})
	if err != nil {
		switch {
		case errors.Is(err, apis.ErrNotFound):
			return err
		case errors.Is(err, apis.ErrRateLimited):
			return err
		default:
			c.Logger().Error(err)
			return err
		}
	}

	if err := storage.AddAddress(session.UserId, models.Address{
		AddressId:   uint8(id),
		Contact:     addressData.Contact,
		CountryCode: addressData.CountryCode,
		Phone:       addressData.Phone,
		Country:     address.Country,
		Street:      address.Street,
		Region:      address.Region,
		City:        address.City,
		Zipcode:     address.Zipcode,
	}, false); err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}

func isCountryCodeValid(code string) bool {
	switch code {
	case "AL", "LT", "LV", "EE", "MD", "RS", "ME", "XK", "BA", "MK", "LI":
		return true
	default:
		return false
	}
}

func checkValue(f url.Values, v string) (string, error) {
	if !f.Has(v) {
		return "", newHTTPError(http.StatusBadRequest, "form has missing '%s' field", v)
	}
	return f.Get(v), nil
}

func getAddressData(c echo.Context) (AddressData, error) {
	form, err := c.FormParams()

	if err != nil {
		return AddressData{}, err
	}

	code, err := checkValue(form, "code")
	if err != nil {
		return AddressData{}, err
	}

	contact, err := checkValue(form, "contact")
	if err != nil {
		return AddressData{}, err
	}

	phone, err := checkValue(form, "phone")
	if err != nil {
		return AddressData{}, err
	}

	street, err := checkValue(form, "street")
	if err != nil {
		return AddressData{}, err
	}

	region, err := checkValue(form, "region")
	if err != nil {
		return AddressData{}, err
	}

	city, err := checkValue(form, "city")
	if err != nil {
		return AddressData{}, err
	}

	zipcode, err := checkValue(form, "zipcode")
	if err != nil {
		return AddressData{}, err
	}

	if !isCountryCodeValid(code) {
		return AddressData{}, newHTTPError(http.StatusBadRequest, "received invalid 'country' field")
	}

	return AddressData{
		CountryCode: code,
		Contact:     contact,
		Phone:       phone,
		Street:      street,
		Region:      region,
		City:        city,
		Zipcode:     zipcode,
	}, nil
}
