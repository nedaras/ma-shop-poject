package handlers

import (
	"nedas/shop/pkg/models"
	"nedas/shop/src/components"
	"net/http"

	"github.com/labstack/echo/v4"
)

func HandleFormAddress(c echo.Context) error {
	id := c.Param("id")

	session := getSession(c)
	storage := getStorage(c)

	if session == nil {
		// todo: idk what todo
		return c.NoContent(http.StatusUnauthorized)
	}

	_ = id
	_ = storage
	return render(c, components.AddressForm(models.Address{
		Contact:     "Kazys Grinius",
		Phone:       "+370 000000",
		CountryCode: "LV",
		Street:      "oblock",
		City:        "Blenkoks",
		Region:      "Chicago",
		Country:     "USA",
		Zipcode:     "54321",
	}))
}
