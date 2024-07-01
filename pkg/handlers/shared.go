package handlers

import (
	"fmt"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

type ErrResult struct {
	Val any
	Err error
}

func render(c echo.Context, comp templ.Component) error {
	return comp.Render(c.Request().Context(), c.Response())
}

func renderWithStatus(sc int, c echo.Context, comp templ.Component) error {
	c.Response().Status = sc
	return render(c, comp)
}

func newHTTPError(code int, format string, a ...any) *echo.HTTPError {
	return echo.NewHTTPError(code, fmt.Sprintf(format, a...))
}
