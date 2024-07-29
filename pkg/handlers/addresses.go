package handlers

import (
	"nedas/shop/src/views"

	"github.com/labstack/echo/v4"
)

func HandleAddresses(c echo.Context) error {
	session := getSession(c)
	storage := getStorage(c)

	if session == nil {
		return unauthorized(c)
	}

	addresses, err := storage.GetAddresses(session.UserId)
	if err != nil {
		return err
	}

	return render(c, views.Addresses(addresses))
}
