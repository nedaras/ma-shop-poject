package handlers

import (
	"nedas/shop/src/views"

	"github.com/labstack/echo/v4"
)

func HandleAddress(c echo.Context) error {
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
			return render(c, views.Address(adress))
		}
	}

	return redirect(c, "/addresses", views.Addresses(addresses))
}
