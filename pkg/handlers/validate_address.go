package handlers

import (
	"errors"
	"fmt"
	"nedas/shop/pkg/apis"
	"nedas/shop/pkg/models"
	"net/http"
	"strings"

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

func HandlePutAddress(c echo.Context) error {
	session := getSession(c)
	storage := getStorage(c)

	if session == nil {
		return unauthorized(c)
	}

	id := c.Param("id")
	addressData, err := validateAddressData(c)
	if err != nil {
		return err
	}

	// todo: like give a check if addr even changed
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
			// todo: yee handle these
			return err
		case errors.Is(err, apis.ErrRateLimited):
			fmt.Println("time till next req:", apis.GetTimeTillNextRequest())
			// todo: yee handle these
			return err
		default:
			c.Logger().Error(err)
			return err
		}
	}

	err = storage.AddAddress(session.UserId, models.Address{
		AddressId:   id,
		Contact:     addressData.Contact,
		CountryCode: addressData.CountryCode,
		Phone:       addressData.Phone,
		Country:     address.Country,
		Street:      address.Street,
		Region:      address.Region,
		City:        address.City,
		Zipcode:     address.Zipcode,
	})
	if err != nil {
		return err
	}

	return redirect(c, "/addresses")
}

func isCountryCodeValid(code string) bool {
	switch code {
	case "AL", "LT", "LV", "EE", "MD", "RS", "ME", "XK", "BA", "MK", "LI":
		return true
	default:
		return false
	}
}

// todo: unit test this validator
func validateAddressData(c echo.Context) (AddressData, error) {
	code := c.FormValue("code")
	if code == "" {
		return AddressData{}, newHTTPError(http.StatusBadRequest, "form has missing 'code' field")
	}

	if !isCountryCodeValid(code) {
		return AddressData{}, newHTTPError(http.StatusBadRequest, "received invalid 'code' field")
	}

	contact := c.FormValue("contact")
	if contact == "" {
		return AddressData{}, newHTTPError(http.StatusBadRequest, "form has missing 'contact' field")
	}

	if len(contact) > 64 {
		return AddressData{}, newHTTPError(http.StatusBadRequest, "received invalid 'contact' field")
	}

	phone := c.FormValue("phone")
	if phone == "" {
		return AddressData{}, newHTTPError(http.StatusBadRequest, "form has missing 'phone' field")
	}

	phone, ok := validatePhone(phone)
	if !ok {
		return AddressData{}, newHTTPError(http.StatusBadRequest, "received invalid 'phone' field")
	}

	street := c.FormValue("street")
	if street == "" {
		return AddressData{}, newHTTPError(http.StatusBadRequest, "form has missing 'street' field")
	}

	region := c.FormValue("region")
	if region == "" {
		return AddressData{}, newHTTPError(http.StatusBadRequest, "form has missing 'region' field")
	}

	city := c.FormValue("city")
	if city == "" {
		return AddressData{}, newHTTPError(http.StatusBadRequest, "form has missing 'city' field")
	}

	zipcode := c.FormValue("zipcode")
	if zipcode == "" {
		return AddressData{}, newHTTPError(http.StatusBadRequest, "form has missing 'zipcode' field")
	}

	if len(zipcode) > 5 || len(zipcode) < 4 {
		return AddressData{}, newHTTPError(http.StatusBadRequest, "received invalid 'zipcode' field")
	}

	for i := range zipcode {
		if zipcode[i] > '9' || zipcode[i] < '0' {
			return AddressData{}, newHTTPError(http.StatusBadRequest, "received invalid 'zipcode' field")
		}
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

func validatePhone(phone string) (string, bool) {
	phone = strings.ReplaceAll(phone, " ", "")
	if phone == "" {
		return "", false
	}

	i := 0
	if phone[0] == '+' {
		i++
	}

	if len(phone[i:]) > 14 || len(phone[i:]) < 7 {
		return "", false
	}

	for i < len(phone) {
		if phone[i] < '0' || phone[i] > '9' {
			return "", false
		}
		i++
	}
	return phone, true
}
