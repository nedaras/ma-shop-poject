package handlers

import (
	"nedas/shop/src/components"
	"net/http"

	"github.com/labstack/echo/v4"
)

func HandlePathValidation(c echo.Context) error {
  url := c.FormValue("url")
  if (url == "") {
    return newHTTPError(http.StatusBadRequest, "field 'url' is emty or not defined");
  }
  return render(c, components.Test(url))

}
