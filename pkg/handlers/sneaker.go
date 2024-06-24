package handlers

import (
	"encoding/json"
	"fmt"
	"nedas/shop/src/components"
	"net/http"

	"github.com/labstack/echo/v4"
)

type SneakerData struct {
  Objects []struct {
    Imagery []struct {
      ViewCode string `json:"viewCode"`
      ViewNumber string `json:"viewNumber"`
      ImageSourceURL string `json:"imageSourceURL"`
      ImageSourceURLTemplate string `json:"imageSourceURLTemplate"`
    } `json:"imagery"`
  } `json:"objects"`
  Errors []interface{} `json:"errors"`
}

func getImageByID(sd *SneakerData, id string) (string, error) {
  for _, i := range sd.Objects {
    for _, image := range i.Imagery {
      if image.ViewCode == id {
        return image.ImageSourceURL, nil
      }
    }
  }
  return "", fmt.Errorf("image not found")
}

func HandleSneaker(c echo.Context) error {
  url := c.FormValue("url")
  if (url == "") {
    return newHTTPError(http.StatusBadRequest, "field 'url' is empty or not defined");
  }

  // and we need to validate domain first
  // idk how to get the id, from the url scrape that shit we will even get the shoe title
  // use regex for that one field with a class

  // we need templates for errors
  res, err := http.Get(fmt.Sprintf("https://api.nike.com/customization/consumer_designs/v1?filter=shortId(%s)", url))
  if err != nil {
    return err
  }

  decoder := json.NewDecoder(res.Body)

  var response SneakerData
  if err := decoder.Decode(&response); err != nil {
    return err
  }

  fmt.Println(response)

  src, err := getImageByID(&response, "B")
  if err != nil {
    // give some templ or shit
    return err
  }

  sc := components.SneakerContext {
    Title: "Some shoe",
    ImageSrc: src,
  }

  return render(c, components.Sneaker(sc))
}
