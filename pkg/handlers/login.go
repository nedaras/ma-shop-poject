package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"nedas/shop/pkg/session"
	"nedas/shop/src/views"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
)

type GoogleAuthData struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	IDToken      string `json:"id_token"`
	ExpiresIn    int    `json:"expires_in"`
}

var (
	ErrInvalidCode = errors.New("provided code is invalid")
	scopes         = "https://www.googleapis.com/auth/userinfo.email+https://www.googleapis.com/auth/userinfo.profile"
)

func HandleLogin(c echo.Context) error {
	session := getSession(c)
	if session == nil {
		return render(c, views.Login(getGoogleLoginURL(scopes)))
	}
	return renderSimpleError(c, http.StatusNotFound)
}

func HandleGoogleLogin(c echo.Context) error {
	code := c.QueryParam("code")
	if code == "" {
		return renderSimpleError(c, http.StatusNotFound)
	}

	data, err := getGoogleAuthData(code)
	if err != nil {
		if errors.Is(err, ErrInvalidCode) {
			return renderSimpleError(c, http.StatusNotFound)
		}
		c.Logger().Error(err)
		return renderSimpleError(c, http.StatusInternalServerError)
	}
	id, err := getGoogleUserID(data)
	if err != nil {
		c.Logger().Error(err)
		return renderSimpleError(c, http.StatusInternalServerError)
	}

	// prob hash the id for better locality
	// save to database
	// generate session token or smth

	session := session.NewSession(id)
	c.SetCookie(session.Cookie())

	return c.Redirect(http.StatusMovedPermanently, "/")
}

// Any returned error will be of type [*OAuth2Error].
func getGoogleUserID(d *GoogleAuthData) (string, error) {
	// we coould decode jwt but idk idk to much work or we can do unsafe way but idk idk
	// for 0 users dont get over my self
	url := "https://www.googleapis.com/oauth2/v1/userinfo?alt=json&access_token=" + d.AccessToken
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", &OAuth2Error{Provider: "GOOGLE", URL: url, Err: err}
	}

	req.Header.Set("Authorization", "Bearer "+d.IDToken)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", &OAuth2Error{Provider: "GOOGLE", URL: url, Err: err}
	}
	defer res.Body.Close()

	data := &struct {
		ID string `json:"id"`
	}{}
	decoder := json.NewDecoder(res.Body)

	if res.StatusCode != 200 {
		switch res.StatusCode {
		case 401:
			return "", &OAuth2Error{Provider: "GOOGLE", URL: url, Err: ErrInvalidCode} // mb remame this error idk
		default:
			return "", &OAuth2Error{Provider: "GOOGLE", URL: url, Err: fmt.Errorf("got unexpected response code '%d'", res.StatusCode)}
		}
	}

	if res.Header.Get("Content-Type") != "application/json; charset=UTF-8" {
		return "", &OAuth2Error{Provider: "GOOGLE", URL: url, Err: errors.New("responded content is not in UTF-8 json form")}
	}

	if err := decoder.Decode(data); err != nil {
		return "", &OAuth2Error{Provider: "GOOGLE", URL: url, Err: err}
	}

	if data.ID == "" {
		return "", &OAuth2Error{Provider: "GOOGLE", URL: url, Err: errors.New("user id is empty")}
	}

	return data.ID, nil
}

// Any returned error will be of type [*OAuth2Error].
func getGoogleAuthData(code string) (*GoogleAuthData, error) {
	if code == "" {
		panic("passed in an empty google code")
	}

	url := fmt.Sprintf("https://oauth2.googleapis.com/token?redirect_uri=%s&client_id=%s&client_secret=%s&code=%s&grant_type=authorization_code",
		os.Getenv("GOOGLE_REDIRECT_URL"),
		os.Getenv("GOOGLE_CLIENT_ID"),
		os.Getenv("GOOGLE_CLIENT_SECRET"),
		code,
	)

	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return nil, &OAuth2Error{Provider: "GOOGLE", URL: url, Err: err}
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, &OAuth2Error{Provider: "GOOGLE", URL: url, Err: err}
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		switch res.StatusCode {
		case 400:
			return nil, &OAuth2Error{Provider: "GOOGLE", URL: url, Err: ErrInvalidCode}
		default:
			return nil, &OAuth2Error{Provider: "GOOGLE", URL: url, Err: fmt.Errorf("got unexpected response code '%d'", res.StatusCode)}
		}
	}

	if res.Header.Get("Content-Type") != "application/json; charset=utf-8" {
		return nil, &OAuth2Error{Provider: "GOOGLE", URL: url, Err: errors.New("responded content is not in UTF-8 json form")}
	}

	data := new(GoogleAuthData)
	decoder := json.NewDecoder(res.Body)

	if err := decoder.Decode(data); err != nil {
		return nil, &OAuth2Error{Provider: "GOOGLE", URL: url, Err: err}
	}

	if data.AccessToken == "" || data.IDToken == "" {
		return nil, &OAuth2Error{Provider: "GOOGLE", URL: url, Err: errors.New("one of the responded fields were empty")}
	}

	return data, nil
}

func getGoogleLoginURL(scopes string) string {
	return fmt.Sprintf("https://accounts.google.com/o/oauth2/v2/auth?redirect_uri=%s&response_type=code&include_granted_scopes=true&access_type=offline&promt=consent&client_id=%s&scope=%s",
		os.Getenv("GOOGLE_REDIRECT_URL"),
		os.Getenv("GOOGLE_CLIENT_ID"),
		scopes,
	)
}

type OAuth2Error struct {
	Provider string
	URL      string
	Err      error
}

func (e *OAuth2Error) Error() string {
	return e.Provider + " '" + e.URL + "': " + e.Err.Error()
}

func (e *OAuth2Error) Unwrap() error {
	return e.Err
}
