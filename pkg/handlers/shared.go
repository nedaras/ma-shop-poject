package handlers

import (
	"fmt"
	"nedas/shop/pkg/middlewares"
	"nedas/shop/src/views"
	"net/http"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

type Auth = middlewares.Auth

type ErrResult[T any] struct {
	Val T
	Err error
}

func getAuth(c echo.Context) Auth {
	val, ok := c.Get("Authentication").(Auth)
	if !ok {
		panic("not using auth middleware")
	}
	return val
}

func render(c echo.Context, comp templ.Component) error {
	return comp.Render(c.Request().Context(), c.Response())
}

func renderWithStatus(sc int, c echo.Context, comp templ.Component) error {
	c.Response().Status = sc
	return render(c, comp)
}

func renderSimpleError(c echo.Context, sc int) error {
	return renderWithStatus(sc, c, views.Error(sc, http.StatusText(sc)))
}

func newHTTPError(code int, format string, a ...any) *echo.HTTPError {
	return echo.NewHTTPError(code, fmt.Sprintf(format, a...))
}
