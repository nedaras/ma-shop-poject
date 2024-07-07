package middlewares

import (
	"nedas/shop/pkg"

	"github.com/labstack/echo/v4"
)

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("auth-session")
		if err != nil {
			return next(c)
		}

		session, ok := session.SessionFromHash(cookie.Value)
		if !ok {
			return next(c)
		}

		c.Set("auth-session", session)

		return next(c)
	}
}
