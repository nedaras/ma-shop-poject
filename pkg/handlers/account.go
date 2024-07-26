package handlers

import (
	"nedas/shop/pkg/models"
	"nedas/shop/src/views"
	"net/http"

	"github.com/labstack/echo/v4"
)

func HandleAccount(c echo.Context) error {
	session := getSession(c)
	storage := getStorage(c)

	if session == nil {
		c.Response().Header().Add("HX-Push-url", "/login")
		return renderWithStatus(http.StatusSeeOther, c, views.Login())
	}

	user, err := storage.GetUser(session.UserId)
	if err != nil {
		return err
	}

	return render(c, views.Account(views.AccountContext{
		User: user,
		Addresses: []models.Address{{
			Contact: "Nedas Pranskunas",
			Phone:   "+370 4635697",
			Street:  "Sarkuvos gatve 10",
			City:    "Kaunas",
			Region:  "Kauno apskritis",
			Country: "Lietuva",
			Zipcode: "12345",
		}, {
			Contact: "Lauras Pranskunas",
			Phone:   "+370 4635697",
			Street:  "Sarkuvos gatve 10",
			City:    "Kaunas",
			Region:  "Kauno apskritis",
			Country: "Lietuva",
			Zipcode: "54321",
		}},
	}))
}
