package middlewares

import (
	"nedas/shop/pkg/session"

	"github.com/labstack/echo/v4"
)

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("auth-session")
		if err != nil {
			return next(c)
		}

		// todo: we can check for user ip address to see if a key is allowed or sum
		//       it will suck for those who uses vpns i guess
		//       we could add like expire date of like a week or a day and like if its expired and we want to update
		//       just check the ip if it dont match log out the user
		session, ok := session.SessionFromHash(cookie.Value)
		if !ok {
			return next(c)
		}

		c.Set("auth-session", session)

		return next(c)
	}
}
