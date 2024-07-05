package handlers

import (
	"nedas/shop/src/views"
	"os"

	"github.com/labstack/echo/v4"
)

// we need some kind of function to handle this idk what half of flag even do
const (
	base   = "https://accounts.google.com/o/oauth2/v2/auth"
	uri    = "?redirect_uri=http://localhost:3000/login/google"
	scopes = "&scope=https://www.googleapis.com/auth/userinfo.email+https://www.googleapis.com/auth/userinfo.profile"
)

func HandleLogin(c echo.Context) error {
	url := base + uri + scopes + "&access_type=offline&response_type=code&promt=consent&client_id=" + os.Getenv("GOOGLE_CLIENT_ID")
	return render(c, views.Login(url))
}
