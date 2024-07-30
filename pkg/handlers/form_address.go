package handlers

import (
	"nedas/shop/pkg/models"
	"nedas/shop/src/components"

	"github.com/labstack/echo/v4"
)

func HandleFormAddress(c echo.Context) error {
	session := getSession(c)
	storage := getStorage(c)

	if session == nil {
		return unauthorized(c)
	}

	id := c.Param("id")
	addresses, err := storage.GetAddresses(session.UserId)
	if err != nil {
		return err
	}

	for _, adress := range addresses {
		if adress.AddressId == id {
			return render(c, components.AddressForm(adress))
		}
	}

	return render(c, components.AddressForm(models.Address{AddressId: id}))
}
