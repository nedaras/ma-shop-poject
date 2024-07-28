package handlers

import (
	"errors"
	"nedas/shop/src/views"
	"net/http"

	"github.com/labstack/echo/v4"
)

func HandleAccount(c echo.Context) error {
	session := getSession(c)
	storage := getStorage(c)

	if session == nil {
		c.Response().Header().Add("HX-Push-url", "/login")
		return renderWithStatus(http.StatusSeeOther, c, views.Login())
	}

	user, err := storage.GetUser(session.UserId)
	if err != nil {
		if errors.Is(err, StorageErrNotFound) {
			// todo: if 404 always remoe the cookie
			c.Response().Header().Add("HX-Push-url", "/login")
			return renderWithStatus(http.StatusSeeOther, c, views.Login())
		}
		return err
	}

	return render(c, views.Account(user))
}
