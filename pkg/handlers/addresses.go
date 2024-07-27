package handlers

import (
	"nedas/shop/pkg/models"
	"nedas/shop/src/views"

	"github.com/labstack/echo/v4"
)

func HandleAddresses(c echo.Context) error {
	return render(c, views.Addresses([]models.Address{{
		Contact:     "Antanas Smetone",
		Phone:       "+372 654652",
		CountryCode: "EE",
		Street:      "oblock",
		City:        "Kaunas",
		Region:      "?",
		Country:     "BBZ",
		Zipcode:     "12345",
	}, {
		Contact:     "Kazys Grinius",
		Phone:       "+370 000000",
		CountryCode: "EE",
		Street:      "oblock",
		City:        "Blenkoks",
		Region:      "Chicago",
		Country:     "USA",
		Zipcode:     "54321",
	}}))

}
