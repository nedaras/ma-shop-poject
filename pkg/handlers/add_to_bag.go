package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

func AddToBag(c echo.Context) error {
	product, err := getQueryProduct(c)
	if err != nil {
		return err
	}

	size, gender := c.FormValue("size"), strings.ToLower(c.FormValue("gender"))
	if size == "" {
		return newHTTPError(http.StatusBadRequest, "form param 'size' is not specified")
	}

	if gender == "" {
		return newHTTPError(http.StatusBadRequest, "form param 'gender' is not specified")
	}

	if gender != "men" && gender != "women" {
		return newHTTPError(http.StatusBadRequest, "query param 'gender' is invalid")
	}

	if len(size) > 4 {
		// would be better like parsing an float to check if its good
		return newHTTPError(http.StatusBadRequest, "query param 'size' is invalid")
	}

	sizes, err := GetSizes(product.PathName, gender == "men")
	if err != nil {
		return err
	}

	for _, s := range sizes {
		if size == s {
			return c.NoContent(http.StatusOK)
		}
	}

	return newHTTPError(http.StatusBadRequest, "query param 'size' is invalid")
}
