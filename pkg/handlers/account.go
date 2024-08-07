package handlers

import (
	"errors"
	"nedas/shop/src/views"

	"github.com/labstack/echo/v4"
)

func HandleAccount(c echo.Context) error {
	session := getSession(c)
	storage := getStorage(c)

	if session == nil {
		return unauthorized(c)
	}

	user, err := storage.GetUser(session.UserId)
	if err != nil {
		if errors.Is(err, StorageErrNotFound) {
			return unauthorized(c)
		}
		return err
	}

	addresses, err := storage.GetAddresses(session.UserId)
	if err != nil {
		return err
	}

	return render(c, views.Account(views.AccountContext{
		User:      user,
		Addresses: addresses,
	}))
}
