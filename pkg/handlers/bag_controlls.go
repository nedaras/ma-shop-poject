package handlers

import (
	"errors"
	"fmt"
	"nedas/shop/pkg/models"
	"nedas/shop/src/components"
	"net/http"

	"github.com/labstack/echo/v4"
)

func HandleProduct(c echo.Context) error {
	session := getSession(c)
	storage := getStorage(c)

	if session == nil {
		return unauthorized(c)
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
		amount, err := storage.AddProduct(session.UserId, product.ThreadId, product.Mid, size)
		if err != nil {
			return err
		}
		return renderWithStatus(http.StatusAccepted, c, components.BagProduct(components.BagProductContext{
			Product:     product,
			Size:        size,
			Amount:      amount,
			RedirectURL: "/" + product.ThreadId + "/" + product.Mid,
		}))
	case http.MethodDelete:
		if err := storage.DeleteProduct(session.UserId, product.ThreadId, product.Mid, size); err != nil {
			return err
		}
		return c.NoContent(http.StatusOK)
	default:
		panic("got unexpected method")
	}
}

func HandleIncrement(c echo.Context) error {
	session := getSession(c)
	storage := getStorage(c)

	if session == nil {
		return unauthorized(c)
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
		if errors.Is(err, StorageErrNotFound) {
			return c.NoContent(http.StatusNotFound)
		}
		return err
	}

	return render(c, components.BagProduct(components.BagProductContext{
		Product:     product,
		Size:        size,
		Amount:      amount,
		RedirectURL: "/" + product.ThreadId + "/" + product.Mid,
	}))
}

func HandleDecrement(c echo.Context) error {
	session := getSession(c)
	storage := getStorage(c)

	if session == nil {
		return unauthorized(c)
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
		if errors.Is(err, StorageErrNotFound) {
			return c.NoContent(http.StatusNotFound)
		}
		return err
	}

	if amount == 0 {
		return c.NoContent(http.StatusOK)
	}

	return render(c, components.BagProduct(components.BagProductContext{
		Product:     product,
		Size:        size,
		Amount:      amount,
		RedirectURL: "/" + product.ThreadId + "/" + product.Mid,
	}))
}

func getQueryProduct(c echo.Context) (models.Product, error) {
	tid, mid := c.QueryParam("tid"), c.QueryParam("mid")
	if tid == "" {
		return models.Product{}, newHTTPError(http.StatusBadRequest, "query param 'tid' is not specified")
	}
	if mid == "" {
		return models.Product{}, newHTTPError(http.StatusBadRequest, "query param 'mid' is not specified")
	}

	product, err := getProduct(tid + ":" + mid)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return models.Product{}, newHTTPError(http.StatusNotFound, fmt.Sprintf("product not found with thread id '%s' and mid '%s'", tid, mid))
		}
		return models.Product{}, err
	}
	return product, nil
}
