package handlers

import (
	"errors"
	"fmt"
	"nedas/shop/src/components"
	"net/http"

	"github.com/labstack/echo/v4"
)

func HandleIncrement(c echo.Context) error {
	session := getSession(c)
	storage := getStorage(c)

	if session == nil {
		// todo: do the cookie stuff
		return newHTTPError(http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
	}

	product, err := getQueryProduct(c)
	if err != nil {
		return err
	}

	amount, err := storage.IncreaseProduct(session.UserId, product.ThreadId, product.Mid)
	if err != nil {
		return err
	}

	return render(c, components.BagProduct(components.BagProductContext{
		Product: product,
		Amount:  amount,
	}))
}

func HandleDecrement(c echo.Context) error {
	session := getSession(c)
	storage := getStorage(c)

	if session == nil {
		// todo: do the cookie stuff
		return newHTTPError(http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
	}

	product, err := getQueryProduct(c)
	if err != nil {
		return err
	}

	amount, err := storage.DecreaseProduct(session.UserId, product.ThreadId, product.Mid)
	if err != nil {
		return err
	}

	if amount == 0 {
		return c.NoContent(http.StatusOK)
	}

	return render(c, components.BagProduct(components.BagProductContext{
		Product: product,
		Amount:  amount,
	}))
}

func HandleDelete(c echo.Context) error {
	session := getSession(c)
	storage := getStorage(c)

	if session == nil {
		return newHTTPError(http.StatusNotImplemented, http.StatusText(http.StatusNotImplemented))
	}

	product, err := getQueryProduct(c)
	if err != nil {
		return err
	}

	if err := storage.DeleteProduct(session.UserId, product.ThreadId, product.Mid); err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}

func getQueryProduct(c echo.Context) (components.Product, error) {
	tid, mid := c.QueryParam("tid"), c.QueryParam("mid")
	if tid == "" {
		return components.Product{}, newHTTPError(http.StatusBadRequest, "query param 'tid' is not specified")
	}
	if mid == "" {
		return components.Product{}, newHTTPError(http.StatusBadRequest, "query param 'mid' is not specified")
	}

	product, err := getProduct(tid + ":" + mid)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return components.Product{}, newHTTPError(http.StatusNotFound, fmt.Sprintf("product not found with thread id '%s' and mid '%s'", tid, mid))
		}
		return components.Product{}, err
	}
	return product, nil
}
