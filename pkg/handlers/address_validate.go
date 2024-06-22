package handlers

import (
	"nedas/shop/src/components"

	"github.com/labstack/echo/v4"
)

// if error return 400 err code and html to update address or sum
func HandleAddressValidate(c echo.Context) error {
  _, err := c.FormParams()
  if err != nil {
    return err
  }

  c.Response().Status = 422;
  return render(c, components.AddressField());

}
