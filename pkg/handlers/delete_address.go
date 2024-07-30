package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func HandleDeleteAddress(c echo.Context) error {
	session := getSession(c)
	storage := getStorage(c)

	if session == nil {
		return unauthorized(c)
	}

	id := c.Param("id")
	if err := storage.DeleteAddress(session.UserId, id); err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}
