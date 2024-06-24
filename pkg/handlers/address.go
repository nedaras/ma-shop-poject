package handlers

import (
	"nedas/shop/src/views"

	"github.com/labstack/echo/v4"
)


func HandleAddress(c echo.Context) error {
  return render(c, views.Address())

}
