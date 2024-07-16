package middlewares

import (
	"nedas/shop/pkg/storage"

	"github.com/labstack/echo/v4"
)

func StorageMiddleware(storage storage.Storage) func(echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("storage", storage)
			return next(c)
		}
	}
}
