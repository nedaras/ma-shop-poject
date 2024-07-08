package handlers

import (
	"fmt"
	"nedas/shop/pkg/session"
	"nedas/shop/pkg/storage"
	"nedas/shop/src/views"
	"net/http"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

type Session = session.Session
type ErrResult[T any] struct {
	Val T
	Err error
}

// optional
func getSession(c echo.Context) *Session {
	val, ok := c.Get("auth-session").(*Session)
	if !ok {
		return nil
	}
	return val
}

func getStorage(c echo.Context) storage.Storage {
	val, ok := c.Get("storage").(storage.Storage)
  if !ok {
    panic("not using storage middleware") 
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
