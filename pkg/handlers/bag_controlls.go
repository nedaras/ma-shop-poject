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

	size := c.QueryParam("size")

	if size == "" {
		return newHTTPError(http.StatusBadRequest, "query param 'size' is not specified")
	}

	product, err := getQueryProduct(c)
	if err != nil {
		return err
	}

	ok, err := validateSize(product.PathName, size)
	if err != nil {
		return err
	}

	if !ok {
		return newHTTPError(http.StatusBadRequest, "query param 'size' is invalid")
	}

	switch c.Request().Method {
	case http.MethodPut:
		if err := storage.AddProduct(session.UserId, product.ThreadId, product.Mid, size); err != nil {
			if errors.Is(err, ErrAlreadySet) {
				amount, err := storage.GetProductAmount(session.UserId, product.ThreadId, product.Mid, size)
				if err != nil {
					return err
				}
				// cuz u know its not rly an error
				return renderWithStatus(http.StatusAccepted, c, components.BagProduct(components.BagProductContext{
					Product: product,
					Size:    size,
					Amount:  amount,
				}))
			}
			return err
		}
	case http.MethodDelete:
		if err := storage.DeleteProduct(session.UserId, product.ThreadId, product.Mid, size); err != nil {
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

	size := c.QueryParam("size")
	if size == "" {
		return newHTTPError(http.StatusBadRequest, "query param 'size' is not specified")
	}

	product, err := getQueryProduct(c)
	if err != nil {
		return err
	}

	ok, err := validateSize(product.PathName, size)
	if err != nil {
		return err
	}

	if !ok {
		return newHTTPError(http.StatusBadRequest, "query param 'size' is invalid")
	}

	amount, err := storage.IncreaseProduct(session.UserId, product.ThreadId, product.Mid, size)
	if err != nil {
		if errors.Is(err, ErrRowNotFound) {
			//return newHTTPError(http.StatusNotFound, fmt.Sprintf("product with thread id '%s' and mid '%s' is not in the bag", product.ThreadId, product.Mid))
			return c.NoContent(http.StatusNotFound)
		}
		return err
	}

	return render(c, components.BagProduct(components.BagProductContext{
		Product: product,
		Size:    size,
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

	size := c.QueryParam("size")
	if size == "" {
		return newHTTPError(http.StatusBadRequest, "query param 'size' is not specified")
	}

	product, err := getQueryProduct(c)
	if err != nil {
		return err
	}

	ok, err := validateSize(product.PathName, size)
	if err != nil {
		return err
	}

	if !ok {
		return newHTTPError(http.StatusBadRequest, "query param 'size' is invalid")
	}

	amount, err := storage.DecreaseProduct(session.UserId, product.ThreadId, product.Mid, size)
	if err != nil {
		if errors.Is(err, ErrRowNotFound) {
			//return newHTTPError(http.StatusNotFound, fmt.Sprintf("product with thread id '%s' and mid '%s' is not in the bag", product.ThreadId, product.Mid))
			return c.NoContent(http.StatusNotFound)
		}
		return err
	}

	if amount == 0 {
		return c.NoContent(http.StatusOK)
	}

	return render(c, components.BagProduct(components.BagProductContext{
		Product: product,
		Size:    size,
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
