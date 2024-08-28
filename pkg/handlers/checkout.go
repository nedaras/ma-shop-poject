package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"nedas/shop/pkg/utils"
	"nedas/shop/src/components"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/stripe/stripe-go/v79"
	"github.com/stripe/stripe-go/v79/checkout/session"
	"github.com/stripe/stripe-go/v79/price"
	"github.com/stripe/stripe-go/v79/product"
)

type ProductBody struct {
	ThreadId string `json:"tid"`
	Mid      string `json:"mid"`
	Amount   string `json:"amount"`
	Size     string `json:"size"`
}

func HandleCheckout(c echo.Context) error {
	session := getSession(c)
	storage := getStorage(c)

	if session == nil {
		return unauthorized(c)
	}

	addressId := c.FormValue("address_id")
	if addressId == "" {
		addresses, err := storage.GetAddresses(session.UserId)
		if err != nil {
			return err
		}
		return render(c, components.AddressSelector(addresses))
	}

	// todo: where the f do i store the address validate it

	products, err := getParamProducts(c.FormValue("products"))
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return newHTTPError(http.StatusNotFound, "product could not be found")
		}
		utils.Logger().Error(err)
		return err
	}

	if len(products) == 0 {
		return newHTTPError(http.StatusNotFound, http.StatusText(http.StatusNotFound))
	}

	var (
		totalPrice float64
	)

	for _, p := range products {
		totalPrice += float64(p.Amount) * p.Product.Price
	}

	url, err := getCheckoutURL(products, "/success?session_id={CHECKOUT_SESSION_ID}&address_id=" + addressId, "/bag")
	if err != nil {
		utils.Logger().Error(err)
		return err
	}

	if isHTMX(c) {
		c.Response().Header().Add("HX-Redirect", url)
		return c.NoContent(http.StatusTemporaryRedirect)
	}

	return c.Redirect(http.StatusTemporaryRedirect, url)
}

func getCheckoutURL(context []components.BagProductContext, success string, cancel string) (string, error) {
	params := make([]*stripe.CheckoutSessionLineItemParams, len(context))
	for i, c := range context {
		product, err := product.New(&stripe.ProductParams{
			Name:        stripe.String(c.Product.Title),
			Description: stripe.String(c.Product.Subtitle + " / UK " + c.Size),
		})

		if err != nil {
			return "", err
		}

		price, err := price.New(&stripe.PriceParams{
			Product: stripe.String(product.ID),
			// todo: make prce like and int cuz floats sucks
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

	a := &stripe.CheckoutSessionParams{
		LineItems:  params,
		SuccessURL: stripe.String(utils.Getenv("HOST") + success),
		CancelURL:  stripe.String(utils.Getenv("HOST") + cancel),
		Mode:       stripe.String(string(stripe.CheckoutSessionModePayment)),
	}

	session, err := session.New(a)
	if err != nil {
		return "", err
	}

	return session.URL, nil
}

func getParamProducts(param string) ([]components.BagProductContext, error) {
	if param == "" {
		return []components.BagProductContext{}, newHTTPError(http.StatusBadRequest, "form has missing 'product' field")
	}

	var products []ProductBody
	if err := json.Unmarshal([]byte(param), &products); err != nil {
		if errors.Is(err, io.EOF) {
			return []components.BagProductContext{}, newHTTPError(http.StatusBadRequest, "'product' field is in invalid json form")
		}
		utils.Logger().Error(err)
		return []components.BagProductContext{}, err
	}

	result := make([]components.BagProductContext, len(products))
	ch := make(chan ErrResult[components.BagProductContext], len(products))

	for _, p := range products {
		if p.ThreadId == "" {
			return []components.BagProductContext{}, newHTTPError(http.StatusBadRequest, "'product' field has missing 'tid' param")
		}
		if p.Mid == "" {
			return []components.BagProductContext{}, newHTTPError(http.StatusBadRequest, "'product' field has missing 'mid' param")
		}

		if p.Amount == "" {
			return []components.BagProductContext{}, newHTTPError(http.StatusBadRequest, "'product' field has missing 'amount' param")
		}

		amount, err := strconv.ParseInt(p.Amount, 10, 8)
		if err != nil {
			return []components.BagProductContext{}, newHTTPError(http.StatusBadRequest, "'product' field has invalid 'amount' param")
		}

		if p.Size == "" {
			return []components.BagProductContext{}, newHTTPError(http.StatusBadRequest, "'product' field has missing 'size' param")
		}
		if len(p.Size) > 4 {
			return []components.BagProductContext{}, newHTTPError(http.StatusBadRequest, "'product' field has invalid 'size' param")
		}
		if _, err := strconv.ParseFloat(p.Amount, 32); err != nil {
			return []components.BagProductContext{}, newHTTPError(http.StatusBadRequest, "'product' field has invalid 'size' param")
		}

		go func() {
			product, err := getProduct(p.ThreadId + ":" + p.Mid)
			if err != nil {
				ch <- ErrResult[components.BagProductContext]{
					Val: components.BagProductContext{},
					Err: err,
				}
			}

			ok, err := validateSize(product.PathName, p.Size)
			if err != nil {
				ch <- ErrResult[components.BagProductContext]{
					Val: components.BagProductContext{},
					Err: err,
				}
			}

			if !ok {
				ch <- ErrResult[components.BagProductContext]{
					Val: components.BagProductContext{},
					Err: ErrNotFound,
				}
			}

			ch <- ErrResult[components.BagProductContext]{
				Val: components.BagProductContext{
					Product: product,
					Size:    p.Size,
					Amount:  uint8(amount),
				},
				Err: nil,
			}

		}()
	}

	for i := range len(products) {
		res := <-ch
		if res.Err != nil {
			return []components.BagProductContext{}, res.Err
		}
		result[i] = res.Val
	}

	return result, nil
}
