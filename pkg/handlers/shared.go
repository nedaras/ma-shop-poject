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

var (
	StorageErrNotFound = storage.ErrNotFound
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

func renderError(c echo.Context, sc int, msg string) error {
	return renderWithStatus(sc, c, views.Error(sc, msg))
}

func renderSimpleError(c echo.Context, sc int) error {
	return renderError(c, sc, http.StatusText(sc))
}

func newHTTPError(code int, format string, a ...any) *echo.HTTPError {
	return echo.NewHTTPError(code, fmt.Sprintf(format, a...))
}

func isHTMX(c echo.Context) bool {
	return c.Request().Header.Get("HX-Request") == "true"
}

func redirect(c echo.Context, path string) error {
	// todo: make rederect render a view i dont like the 2 trafic proxy like
	if isHTMX(c) {
		c.Response().Header().Add("HX-Location", path)
		return c.NoContent(http.StatusSeeOther)
	}
	return c.Redirect(http.StatusSeeOther, path)
}

func redirectB(c echo.Context, path string, comp templ.Component) error {
	if isHTMX(c) {
		c.Response().Header().Add("HX-Push-Url", path)
		c.Response().Header().Add("HX-Retarget", "body")
		c.Response().Header().Add("HX-Reswap", "innerHTML")
		return render(c, comp)
	}
	return c.Redirect(http.StatusSeeOther, path)
}

func unauthorized(c echo.Context) error {
	if isHTMX(c) {
		// todo: i dont like the redirect when we can manipulate boost with headers
		//c.Response().Header().Add("HX-Location", "{\"path\":\"/login\",\"values\":{\"fallback\":\""+fallback+"\"}}")
		c.Response().Header().Add("HX-Location", "/login")
		return c.NoContent(http.StatusUnauthorized)
	}
	return c.Redirect(http.StatusSeeOther, "/login")
}
