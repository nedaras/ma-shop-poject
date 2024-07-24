package handlers

import (
	"nedas/shop/src/components"
	"net/http"

	"github.com/labstack/echo/v4"
)

func AddToBag(c echo.Context) error {
	session := getSession(c)
	storage := getStorage(c)

	if session == nil {
    // todo: move to login field then
		return newHTTPError(http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
	}

	product, err := getQueryProduct(c)
	if err != nil {
		return err
	}

	size := c.FormValue("size")
	if size == "" {
		return newHTTPError(http.StatusBadRequest, "form param 'size' is not specified")
	}

	ok, err := validateSize(product.PathName, size)
	if err != nil {
		return err
	}

	if !ok {
		return newHTTPError(http.StatusBadRequest, "query param 'size' is invalid")
	}

	amount, err := storage.AddProduct(session.UserId, product.ThreadId, product.Mid, size)
	if err != nil {
		return err
	}

	return render(c, components.BagProduct(components.BagProductContext{
		Product: product,
		Size:    size,
		Amount:  amount,
	}))
}

func validateSize(path string, size string) (bool, error) {
	if len(size) > 4 {
		return false, nil
	}

	sizes, err := GetAllSizes(path)
	if err != nil {
		return false, err
	}

	for _, s := range sizes {
		if s == size {
			return true, nil
		}
	}

	return false, nil
}
