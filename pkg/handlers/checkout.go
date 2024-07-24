package handlers

import (
	"nedas/shop/src/views"
	"net/http"

	"github.com/labstack/echo/v4"
)

func HandleCheckout(c echo.Context) error {
	session := getSession(c)

	if session == nil {
		// todo: hacker or out of sync reupdate state
		return newHTTPError(http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
	}

	// todo: note we have to make like a cache for what we tryna checkout cuz we dont want any out of sync problems

	return render(c, views.Address())
}
