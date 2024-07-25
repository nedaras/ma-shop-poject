package handlers

import (
	"fmt"
	"nedas/shop/pkg/apis"
	"nedas/shop/src/views"
	"net/http"
	"net/url"

	"github.com/labstack/echo/v4"
)

type AddressData struct {
	Country string
	Contact string
	Phone   string
	Street  string
	Region  string
	City    string
	Zipcode string
}

// pattern ^[A-Za-zÄÖÜäöüßĄČĘĖĮŠŲŪŽąčęėįšųūž ]+$
// if error return 400 err code and html to update addressData or sum
func HandleAddressValidate(c echo.Context) error {
	addressData, err := getAddressData(c)
	if err != nil {
		return err
	}
	// we need to sanatize and validate the shit out of this cuz what the user writes here will go to and database
	// would be crazy if an user writed in like 1k long names or idk phone number without numbers

	// todo: idk how we need to use an interface no?
	adddress, err := apis.ValidateAddress(apis.Address{
		Country: addressData.Country,
		Street:  addressData.Street,
		Region:  addressData.Region,
		City:    addressData.City,
		Zipcode: addressData.Zipcode,
	})
	if err != nil {
		return err
	}

	fmt.Println("adress issss:", adddress.String())

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

	countryCode, ok := getCountryCode(country)
	if !ok {
		return AddressData{}, newHTTPError(http.StatusBadRequest, "received invalid 'country' field")
	}

	return AddressData{
		Country: country,
		Contact: contact,
		Phone:   countryCode + phone,
		Street:  street,
		Region:  region,
		City:    city,
		Zipcode: zipcode,
	}, nil
}
