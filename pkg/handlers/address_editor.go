package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"nedas/shop/pkg/models"
	"nedas/shop/src/components"

	"github.com/labstack/echo/v4"
)

func HandleAddressEditor(c echo.Context) error {
	session := getSession(c)
	storage := getStorage(c)

	if session == nil {
		return unauthorized(c)
	}

	id := c.Param("id")
	if id == "" {
		return render(c, components.AddressEditor(models.Address{AddressId: generateRandomId()}))
	}

	addresses, err := storage.GetAddresses(session.UserId)
	if err != nil {
		return err
	}

	for _, adress := range addresses {
		if adress.AddressId == id {
			return render(c, components.AddressEditor(adress))
		}
	}

	// todo: if pressed edit on non existing address just delete it from the view
	return render(c, components.AddressEditor(models.Address{AddressId: id}))
}

func generateRandomId() string {
	bytes := make([]byte, 8)
	rand.Read(bytes) // will never err

	return hex.EncodeToString(bytes)
}
