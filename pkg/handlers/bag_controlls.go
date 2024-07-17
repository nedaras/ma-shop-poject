package handlers

import (
	"errors"
	"fmt"
	"nedas/shop/pkg/storage"
	"nedas/shop/src/components"
	"net/http"

	"github.com/labstack/echo/v4"
)

var (
	ErrRowNotFound = storage.ErrNotFound
	ErrAlreadySet  = storage.ErrAlreadySet
)

func HandleProduct(c echo.Context) error {
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

	switch c.Request().Method {
	case http.MethodPut:
		if err := storage.AddProduct(session.UserId, product.ThreadId, product.Mid); err != nil {
			if errors.Is(err, ErrAlreadySet) {
				return newHTTPError(http.StatusConflict, fmt.Sprintf("product with thread id '%s' and mid '%s' is already in the bag", product.ThreadId, product.Mid))
			}
			return err
		}
	case http.MethodDelete:
		if err := storage.DeleteProduct(session.UserId, product.ThreadId, product.Mid); err != nil {
			return err
		}
	default:
		panic("got unexpected method")
	}

	return c.NoContent(http.StatusOK)

}

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
		if errors.Is(err, ErrRowNotFound) {
			return newHTTPError(http.StatusNotFound, fmt.Sprintf("product with thread id '%s' and mid '%s' is not in the bag", product.ThreadId, product.Mid))
		}
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
		if errors.Is(err, ErrRowNotFound) {
			return newHTTPError(http.StatusNotFound, fmt.Sprintf("product with thread id '%s' and mid '%s' is not in the bag", product.ThreadId, product.Mid))
		}
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
