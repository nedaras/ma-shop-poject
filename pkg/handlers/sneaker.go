package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"nedas/shop/src/components"
	"net/http"
	"regexp"
	"strings"

	"github.com/labstack/echo/v4"
)

type NikeConsumerData struct {
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

type NikeScrapedData struct {
  ID string
  Title string
}


func HandleSneaker(c echo.Context) error {
  url := c.FormValue("url")
  if !isNikeURL(url) {
    return newHTTPError(http.StatusBadRequest, "'url' is invalid");
  }

  sData, err := scrapeNikeURL(url);

  if err != nil {
    return err // templ
  }

  cData, err := getNikeConsumerData(sData.ID)
  if err != nil {
    return err // templ
  }

  src, err := getImageByID(cData, "B")
  if err != nil {
    return err // templ
  }

  sc := components.SneakerContext {
    Title: sData.Title,
    ImageSrc: src,
  }

  return render(c, components.Sneaker(sc))
}

func getImageByID(d *NikeConsumerData, id string) (string, error) {
  for _, i := range d.Objects {
    for _, image := range i.Imagery {
      if image.ViewCode == id {
        return image.ImageSourceURL, nil
      }
    }
  }
  return "", fmt.Errorf("image not found")
}

func isNikeURL(url string) bool {
  if strings.HasPrefix(url, "https://www.nike.com/") {
    return true
  }
  if strings.HasPrefix(url, "www.nike.com/") {
    return true
  }
  if strings.HasPrefix(url, "https://nike.com/") {
    return true
  }
  if strings.HasPrefix(url, "nike.com/") {
    return true
  }
  return false
}

func getNikeConsumerData(id string) (*NikeConsumerData, error) {
  res, err := http.Get(fmt.Sprintf("https://api.nike.com/customization/consumer_designs/v1?filter=shortId(%s)", id))
  if err != nil {
    return nil, err
  }
  defer res.Body.Close()

  data := &NikeConsumerData{};
  decoder := json.NewDecoder(res.Body)

  if err := decoder.Decode(&data); err != nil {
    return nil, err
  }
  return data, nil
}

// from scrape we can get everything like prices and shit, but it slow
func scrapeNikeURL(url string) (NikeScrapedData, error) {
  titleR, err := regexp.Compile(`<h1 .*data-test="product-title">(.+)<\/h1>`)
  if err != nil {
    panic("could not compile regexp")
  }

  idR, err := regexp.Compile(`"metricId":"(.{10})"`)
  if err != nil {
    panic("could not compile regexp")
  }

  res, err := http.Get(url)
  if err != nil {
    return NikeScrapedData{}, err
  }
  defer res.Body.Close()

  buf := new(strings.Builder)
  if _, err = io.Copy(buf, res.Body); err != nil {
    return NikeScrapedData{}, err
  }

  titleMatches := titleR.FindStringSubmatch(buf.String())
  if len(titleMatches) != 2 {
    return NikeScrapedData{}, fmt.Errorf("'url' is invalid")
  }

  idMatches := idR.FindStringSubmatch(buf.String())
  if len(idMatches) != 2 {
    return NikeScrapedData{}, fmt.Errorf("'url' is invalid")
  }
  
  return NikeScrapedData {
    Title: titleMatches[1],
    ID: idMatches[1],
  }, nil

}
