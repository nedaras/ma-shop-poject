package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func HandleCheckout(c echo.Context) error {
	session := getSession(c)

	if session == nil {
		return unauthorized(c)
	}

	// todo: note we have to make like a cache for what we tryna checkout cuz we dont want any out of sync problems

	return c.NoContent(http.StatusNotFound)
}
