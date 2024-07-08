package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)
// todo: do all in one not 3 seperate functions for add remove and delete
func HandleIncrement(c echo.Context) error {
  session := getSession(c)
  storage := getStorage(c)

  if session == nil {
    // todo: do the cookie stuff
    return newHTTPError(http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
  }

  // todo idk mb validate if session id is even valid
  // todo add validate id function, it would be faster
  tid, mid := c.QueryParam("tid"), c.QueryParam("mid")
  _, err := getProduct(tid + ":" + mid)

  if err != nil {
    return err
  }

  if err := storage.IncreaseProduct(session.UserId, tid, mid); err != nil {
    return err
  }

  return newHTTPError(http.StatusExpectationFailed, "OK OK OK OK")

}
