package handlers

import (
	"errors"
	"fmt"
	"nedas/shop/pkg/models"
	"nedas/shop/pkg/utils"
	"nedas/shop/src/components"
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

	url, err := getCheckoutURL(products)
	if err != nil {
		utils.Logger().Error(err)
		return err
	}

	return c.Redirect(http.StatusMovedPermanently, url)
}

func getCheckoutURL(context []components.BagProductContext) (string, error) {
	params := make([]*stripe.CheckoutSessionLineItemParams, len(context))
	for i, c := range context {
		product, err := product.New(&stripe.ProductParams{
			Name:        stripe.String(c.Product.Title),
			Description: stripe.String(c.Product.Subtitle),
		})

		if err != nil {
			return "", err
		}

		price, err := price.New(&stripe.PriceParams{
			Product:    stripe.String(product.ID),
			UnitAmount: stripe.Int64(int64(c.Product.Price * 100)),
			Currency:   stripe.String(string(stripe.CurrencyEUR)),
		})

		if err != nil {
			return "", err
		}

		params[i] = &stripe.CheckoutSessionLineItemParams{
			Price:    stripe.String(price.ID),
			Quantity: stripe.Int64(int64(c.Amount)),
		}

	}

	domain := "http://localhost:3000"
	a := &stripe.CheckoutSessionParams{
		LineItems:  params,
		SuccessURL: stripe.String(domain + "/success.html"),
		CancelURL:  stripe.String(domain + "/cancel.html"),
		Mode:       stripe.String(string(stripe.CheckoutSessionModePayment)),
	}

	session, err := session.New(a)
	if err != nil {
		return "", err
	}

	return session.URL, nil
}
