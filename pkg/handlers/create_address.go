package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"nedas/shop/pkg/models"
	"nedas/shop/src/views"

	"github.com/labstack/echo/v4"
)

func HandleCreateAddress(c echo.Context) error {
	session := getSession(c)
	if session == nil {
		return unauthorized(c)
	}

	id := generateRandomId()

	c.Response().Header().Add("HX-Push-Url", "/address/"+id)
	return render(c, views.Address(models.Address{AddressId: id}))
}

func generateRandomId() string {
	bytes := make([]byte, 8)
	rand.Read(bytes)

	return hex.EncodeToString(bytes)
}
