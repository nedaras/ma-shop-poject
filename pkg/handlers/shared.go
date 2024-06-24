package handlers

import (
	"fmt"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

func render(c echo.Context, comp templ.Component) error {
  return comp.Render(c.Request().Context(), c.Response());
}

func newHTTPError(code int, format string, a ...any) *echo.HTTPError {
  return echo.NewHTTPError(code, fmt.Sprintf(format, a...))
}
