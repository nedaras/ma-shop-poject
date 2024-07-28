package handlers

import (
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
		// todo: idk what todo
		return c.NoContent(http.StatusUnauthorized)
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 8)
	if err != nil {
		return newHTTPError(http.StatusBadRequest, "param 'id' is not valid uint8")
	}

	user, err := storage.GetUser(session.UserId)
	if err != nil {
		// todo: handle not found
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
