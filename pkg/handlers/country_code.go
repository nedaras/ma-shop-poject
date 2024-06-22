package handlers

import (
	"fmt"
	"nedas/shop/src/components"

	"github.com/labstack/echo/v4"
)

func getCountryCode(c string) (string, bool) {
  switch c {
    case "": return "", true; // u can never know mb that disabled attr will not work in some browsers
    case "LT": return "+370", true;
    case "LV": return "+371", true;
    case "EE": return "+372", true;
    default: return "", false;
  }
}

// we will only handle errors if it is needed, below errors should never apear if user uses interface as intended
func HandleCountryCode(c echo.Context) error {
  params, err := c.FormParams()
  if err != nil {
    return err
  }

  values, ok := params["country"];
  if !ok {
    return fmt.Errorf("invalid form, missing 'country' tag")
  }

  if len(values) != 1 {
    return fmt.Errorf("invalid form, 'country' tag is corrupted")
  }

  country := values[0]
  code, ok := getCountryCode(country)

  if !ok {
    return fmt.Errorf("country not valid: '%s'", country)
  }
  return render(c, components.CountryCode(code));
}
