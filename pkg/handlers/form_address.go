package handlers

import (
	"errors"
	"nedas/shop/pkg/models"
	"nedas/shop/src/components"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func HandleFormAddress(c echo.Context) error {
	session := getSession(c)
	storage := getStorage(c)

	if session == nil {
		return unauthorized(c)
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 8)
	if err != nil {
		return newHTTPError(http.StatusBadRequest, "param 'id' is not valid uint8")
	}

	user, err := storage.GetUser(session.UserId)
	if err != nil {
		if errors.Is(err, StorageErrNotFound) {
			return unauthorized(c)
		}
		return err
	}

	if len(user.Addresses) == 0 {
		return render(c, components.AddressForm(models.Address{AddressId: uint8(id)}))
	}

	for _, adress := range user.Addresses {
		if adress.AddressId == uint8(id) {
			return render(c, components.AddressForm(adress))
		}
	}

	return newHTTPError(http.StatusNotFound, http.StatusText(http.StatusNotFound))
}
