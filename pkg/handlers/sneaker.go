package handlers

import (
	"errors"
	"nedas/shop/src/views"
	"net/http"

	"github.com/labstack/echo/v4"
)

// https://api.nike.com/cic/grand/v1/graphql/getfulfillmenttypesofferings/v4?variables=%7B%22countryCode%22%3A%22GB%22%2C%22currency%22%3A%22GBP%22%2C%22locale%22%3A%22en-GB%22%2C%22locationId%22%3A%22%22%2C%22locationType%22%3A%22STORE_VIEWS%22%2C%22offeringTypes%22%3A%5B%22SHIP%22%5D%2C%22postalCode%22%3A%22%22%2C%22productId%22%3A%2210c70f8d-07e3-5653-b02c-bae0e5671a45%22%7D
func HandleSneaker(c echo.Context) error {
	tid := c.Param("thread_id")
	mid := c.Param("mid")

	product, err := getProduct(tid + ":" + mid)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return renderSimpleError(c, http.StatusNotFound)
		}
		// todo: all them loggers where idk whats the error
		c.Logger().Error(err)
		return renderSimpleError(c, http.StatusInternalServerError)
	}

	sizes, err := GetSizes(product.PathName, true)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return renderSimpleError(c, http.StatusNotFound)
		}
		c.Logger().Error(err)
		return renderSimpleError(c, http.StatusInternalServerError)
	}

	return render(c, views.Sneaker(views.SneakerContext{
		Product:  product,
		Sizes:    sizes,
		LoggedIn: getSession(c) != nil,
	}))
}
