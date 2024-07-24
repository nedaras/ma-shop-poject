package handlers

import (
	"fmt"
	"nedas/shop/src/views"
	"net/http"
	"net/url"

	"github.com/labstack/echo/v4"
)

type AddressData struct {
	Country  string
	Contact  string
	Phone    string
	Address1 string
	Address2 string
	Region   string
	City     string
	Zipcode  string
}

// if error return 400 err code and html to update addressData or sum
func HandleAddressValidate(c echo.Context) error {
	addressData, err := getAddressData(c)
	if err != nil {
		return err
	}

	fmt.Println(addressData)

	// todo: like i have phone numbers saved i need to save provinces too
	// fuck google i cant add address validation api cuz my card is prepaid or sum f them
	return render(c, views.Index())

}

func getCountryCode(country string) (string, bool) {
	switch country {
	case "AL":
		return "+355", true
	case "LT":
		return "+370", true
	case "LV":
		return "+371", true
	case "EE":
		return "+372", true
	case "MD":
		return "+373", true
	case "RS":
		return "+381", true
	case "ME":
		return "+382", true
	case "XK":
		return "+383", true
	case "BA":
		return "+387", true
	case "MK":
		return "+389", true
	case "LI":
		return "+423", true
	default:
		return "", false
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

	country, err := checkValue(form, "country")
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

	address1, err := checkValue(form, "address_1")
	if err != nil {
		return AddressData{}, err
	}

	address2, err := checkValue(form, "address_2")
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

	countryCode, ok := getCountryCode(country)
	if !ok {
		return AddressData{}, newHTTPError(http.StatusBadRequest, "received invalid 'country' field")
	}

	return AddressData{
		Country:  country,
		Contact:  contact,
		Phone:    countryCode + phone,
		Address1: address1,
		Address2: address2,
		Region:   region,
		City:     city,
		Zipcode:  zipcode,
	}, nil
}
