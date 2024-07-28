package handlers

import (
	"nedas/shop/src/views"
	"net/http"

	"github.com/labstack/echo/v4"
)

func HandleAddresses(c echo.Context) error {
	session := getSession(c)
	storage := getStorage(c)

	// todo: handle out of sync idk
	if session == nil {
		return newHTTPError(http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
	}

	user, err := storage.GetUser(session.UserId)
	if err != nil {
		return err
	}

	return render(c, views.Addresses(user.Addresses))
}
