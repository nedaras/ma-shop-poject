package middlewares

import "github.com/labstack/echo/v4"

type Auth struct {
	Token    string
	LoggedIn bool
}

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Set("Authentication", Auth{
			Token:    "BBD",
			LoggedIn: false,
		})
		return next(c)
	}
}
