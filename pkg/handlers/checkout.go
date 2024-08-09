package handlers

import (
	"errors"
	"fmt"
	"nedas/shop/pkg/models"
	"nedas/shop/pkg/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

func HandleCheckout(c echo.Context) error {
	session := getSession(c)
	storage := getStorage(c)

	if session == nil {
		return unauthorized(c)
	}

	var (
		address    models.Address
		totalPrice float64
	)

	addressId := c.FormValue("address_id")
	if addressId == "" {
		addresses, err := storage.GetAddresses(session.UserId)
		if err != nil {
			utils.Logger().Error(err)
			return err
		}

		if len(addresses) == 0 {
			// todo:
			panic("redirect to add address")
		}

		if len(addresses) > 1 {
			// todo:
			panic("not implemented.")
		}

		address = addresses[0]
	} else {
		a, err := storage.GetAddress(session.UserId, addressId)
		if err != nil {
			if errors.Is(err, StorageErrNotFound) {
				return newHTTPError(http.StatusNotFound, "form has invalid 'address_id'")
			}
			utils.Logger().Error(err)
			return err
		}
		address = a
	}

	fmt.Println(address)

	// todo: when doing checkout we should kinda trust what client sent us what they wanna buy,
	//       we just will need to validate if the purhase is valid
	products, err := getProducts(session.UserId, storage)
	if err != nil {
		utils.Logger().Error(err)
		return err
	}

	if len(products) == 0 {
		return newHTTPError(http.StatusNotFound, http.StatusText(http.StatusNotFound))
	}

	for _, p := range products {
		totalPrice += float64(p.Amount) * p.Product.Price
	}

	fmt.Println("total price", totalPrice)

	return c.NoContent(http.StatusNotFound)
}
