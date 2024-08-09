package handlers

import (
	"errors"
	"fmt"
	"nedas/shop/pkg/models"
	"nedas/shop/pkg/utils"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/stripe/stripe-go/v79"
	"github.com/stripe/stripe-go/v79/checkout/session"
	"github.com/stripe/stripe-go/v79/price"
	"github.com/stripe/stripe-go/v79/product"
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

	secret, err := getClientSecret()
	if err != nil {
		utils.Logger().Error(err)
		return err
	}

	fmt.Println("total price", totalPrice)
	fmt.Println("secret", secret)

	//return c.NoContent(http.StatusNotFound)
	return c.HTML(http.StatusOK, secret)
}

func getClientSecret() (string, error) {
	product, err := product.New(&stripe.ProductParams{Name: stripe.String("T-shirt")})
	if err != nil {
		return "", err
	}

	price, err := price.New(&stripe.PriceParams{
		Product:    stripe.String(product.ID),
		UnitAmount: stripe.Int64(2000),
		Currency:   stripe.String(string(stripe.CurrencyEUR)),
	})

	if err != nil {
		return "", err
	}

	params := &stripe.CheckoutSessionParams{
		UIMode:    stripe.String("embedded"),
		ReturnURL: stripe.String("http://localhost:3000/stripe?session_id={CHECKOUT_SESSION_ID}"),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Price:    stripe.String(price.ID),
				Quantity: stripe.Int64(1),
			},
		},
		Mode: stripe.String(string(stripe.CheckoutSessionModePayment)),
	}

	session, err := session.New(params)
	if err != nil {
		return "", err
	}

	return session.ClientSecret, nil
}
