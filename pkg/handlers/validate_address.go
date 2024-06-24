package handlers

import (
	"nedas/shop/src/views"
	"net/http"
	"net/url"

	"github.com/labstack/echo/v4"
)

type AddressData struct {
  Country string
  Contact string
  Phone string
  Address1 string
  Address2 string
  Region string
  City string
  Zipcode string
}

// if error return 400 err code and html to update addressData or sum
func HandleAddressValidate(c echo.Context) error {
  addressData, err := getAddressData(c)
  if err != nil {
    return err
  }

  // i cant make cloud account for now
  _ = addressData
  return render(c, views.Index());

}

func getCountryCode(country string) (string, bool) {
  switch country {
    case "AL": return "+355", true
    case "LT": return "+370", true
    case "LV": return "+371", true
    case "EE": return "+372", true
    case "MD": return "+373", true
    case "RS": return "+381", true
    case "ME": return "+382", true
    case "XK": return "+383", true
    case "BA": return "+387", true
    case "MK": return "+389", true
    case "LI": return "+423", true
    default: return "", false
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
  addressData := AddressData{}

  if err != nil {
    return addressData, err
  }

  country, err := checkValue(form, "country")
  if (err != nil) {
    return addressData, err
  }

  contact, err := checkValue(form, "contact")
  if (err != nil) {
    return addressData, err
  }

  phone, err := checkValue(form, "phone")
  if (err != nil) {
    return addressData, err
  }

  address1, err := checkValue(form, "address_1")
  if (err != nil) {
    return addressData, err
  }

  address2, err := checkValue(form, "address_2")
  if (err != nil) {
    return addressData, err
  }

  region, err := checkValue(form, "region")
  if (err != nil) {
    return addressData, err
  }

  city, err := checkValue(form, "city")
  if (err != nil) {
    return addressData, err
  }

  zipcode, err := checkValue(form, "zipcode")
  if (err != nil) {
    return addressData, err
  }

  countryCode, ok := getCountryCode(country)
  if (!ok) {
    return addressData, newHTTPError(http.StatusBadRequest, "received invalid 'country' field");
  }

  addressData.Country = country
  addressData.Contact = contact
  addressData.Phone = countryCode + phone 
  addressData.Address1 = address1
  addressData.Address2 = address2
  addressData.Region = region 
  addressData.City = city 
  addressData.Zipcode = zipcode

  return addressData, nil
}

