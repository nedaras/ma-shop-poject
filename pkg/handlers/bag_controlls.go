package handlers

import (
	"errors"
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

	// todo idk mb validate if session id is even valid
	// todo add validate id function, it would be faster
	tid, mid := c.QueryParam("tid"), c.QueryParam("mid")
	product, err := getProduct(tid + ":" + mid)

	if err != nil {
		return err
	}

	amount, err := storage.IncreaseProduct(session.UserId, tid, mid)
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

	// todo idk mb validate if session id is even valid
	// todo add validate id function, it would be faster
	tid, mid := c.QueryParam("tid"), c.QueryParam("mid")
	product, err := getProduct(tid + ":" + mid)

	if err != nil {
		return err
	}

	amount, err := storage.DecreaseProduct(session.UserId, tid, mid)
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

	tid, mid, err := validateAndGetProductID(c)
	if err != nil {
		return err
	}

	if err := storage.DeleteProduct(session.UserId, tid, mid); err != nil {
		return err

	}

	return c.NoContent(http.StatusOK)
}

func validateAndGetProductID(c echo.Context) (string, string, error) {
	tid, mid := c.QueryParam("tid"), c.QueryParam("mid")
	if tid == "" || mid == "" {
		return "", "", newHTTPError(http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
	}

	if _, err := getProduct(tid + ":" + mid); err != nil {
		if errors.Is(err, ErrNotFound) {
			return "", "", newHTTPError(http.StatusNotFound, http.StatusText(http.StatusNotFound))
		}
		return "", "", err
	}

	return tid, mid, nil
}
