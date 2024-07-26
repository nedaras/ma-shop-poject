package handlers

import (
	"nedas/shop/pkg/models"
	"nedas/shop/src/views"

	"github.com/labstack/echo/v4"
)

func HandleAddresses(c echo.Context) error {
	return render(c, views.Addresses([]models.Address{{
		Contact: "Nedas Pranskunas",
		Phone:   "+370 4635697",
		Street:  "Sarkuvos gatve 10",
		City:    "Kaunas",
		Region:  "Kauno apskritis",
		Country: "Lietuva",
		Zipcode: "12345",
	}, {
		Contact: "Lauras Pranskunas",
		Phone:   "+370 4635697",
		Street:  "Sarkuvos gatve 10",
		City:    "Kaunas",
		Region:  "Kauno apskritis",
		Country: "Lietuva",
		Zipcode: "54321",
	}}))

}
