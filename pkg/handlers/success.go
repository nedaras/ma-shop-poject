package handlers

import (
	"nedas/shop/src/views"
	"net/http"

	"github.com/labstack/echo/v4"
)

func HandleSuccess(c echo.Context) error {
  // how to validate if its stripe request or sum
  // if unauthorized and session_id is valid try to login and do some stuff
  addressId := c.QueryParam("address_id")
  if addressId == "" {
    return newHTTPError(http.StatusNotFound, http.StatusText(http.StatusNotFound))
  }

  return render(c, views.Success())
}
