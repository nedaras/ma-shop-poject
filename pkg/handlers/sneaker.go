package handlers

import (
	"bytes"
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

type NextDataData struct {
  Props struct {
    PageProps struct {
      InitialState struct {
        Threads struct {
          Products map[string]ProductData `json:"products"`
        } `json:"Threads"`
      } `json:"initialState"`
    } `json:"pageProps"`
  } `json:"props"`
  Query struct {
    PBID string `json:"pbid"`
  } `json:"query"`
}

type ProductData struct {
  Title string `json:"title"`
  CurrentPrice float64 `json:"currentPrice"`
  PathName string `json:"pathName"`
}

type NikeScrapedData struct {
  ID string
  Title string
  Price float64
}

// need to unit test
func convertLink(str string) (string, error) {
  if len(str) == 0 {
    return "", fmt.Errorf("'url' is invalid")
  }

  i := 0;
  strBuilder := new(strings.Builder)

  switch str[0] {
    case 'h': {
      {
        flag := "http"
        if i + len(flag) > len(str) {
          return "", fmt.Errorf("'url' is invalid")
        }

        if str[i:i+len(flag)] != flag {
          return "", fmt.Errorf("'url' is invalid")
        }
        i += len(flag)

        if i + 1 > len(str) {
          return "", fmt.Errorf("'url' is invalid")
        }

        if str[i] == 's' {
          i++
        }
      }
      {
        flag := "://"
        if i + len(flag) > len(str) {
          return "", fmt.Errorf("'url' is invalid")
        }

        if str[i:i+len(flag)] != flag {
          return "", fmt.Errorf("'url' is invalid")
        }
        i += len(flag)
      }
    }
    fallthrough
    case 'w': {
      flag := "www."
      if i + len(flag) > len(str) {
          return "", fmt.Errorf("'url' is invalid")
      }

      if str[i:i+len(flag)] == flag {
        i += len(flag)
      } else if i == 0 {
          return "", fmt.Errorf("'url' is invalid")
      }
    }
    fallthrough
    case 'n': {
      flag := "nike.com/"
      if i + len(flag) > len(str) {
          return "", fmt.Errorf("'url' is invalid")
      }

      if str[i:i+len(flag)] != flag {
        return "", fmt.Errorf("'url' is invalid")
      }
      i += len(flag)
    }
  }

  if _, err := strBuilder.WriteString("https://www.nike.com/"); err != nil {
    return "", err
  }

  if i + 2 > len(str) {
    return "", fmt.Errorf("'url' is invalid")
  }

  if str[i:i+2] == "u/" {
    goto final 
  }

  for i < len(str) {
    if str[i] == '/' {
      i++;
      break
    }
    i++;
  }

  if i + 2 > len(str) {
    return "", fmt.Errorf("'url' is invalid")
  }

  if str[i:i+2] != "u/" {
    return "", fmt.Errorf("'url' is invalid")
  }

final:
  i += 2
  if _, err := strBuilder.WriteString("gb/u/"); err != nil {
    return "", err
  }

  if _, err := strBuilder.WriteString(str[i:]); err != nil {
    return "", err
  }

  return strBuilder.String(), nil
}

// for some size guide
// https://api.nike.com/customization/availabilities/v1?filter=pathName(af1mid365ho22)&filter=countryCode(GB)&language=en-GB

// here price if found though no Title, we still need to scrape nike itself
// https://api.nike.com/customization/builderaggregator/v2/builder/GB/en_GB/af1mid365ho22

// mb will be usefull oneday
// https://www.nike.com/assets/nikeid/builder-helper/dist/language-mapper/languageMap.json

// there has to be away to get title cuz recomendation in that page

// there is an idea to handle this scrape and stuff in like other languege, like zig
// we need to convert to GB link and in bg check conversion rates
func HandleSneaker(c echo.Context) error {
  url, err := convertLink(c.FormValue("url"))
  if err != nil {
    return err // templ
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
    Price: sData.Price,
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
  nextDataR := regexp.MustCompile(`<script id="__NEXT_DATA__" type="application\/json">(.+)<\/script>`)

  res, err := http.Get(url)
  if err != nil {
    return NikeScrapedData{}, err
  }
  defer res.Body.Close()

  buf := new(strings.Builder)
  if _, err = io.Copy(buf, res.Body); err != nil {
    return NikeScrapedData{}, err
  }

  nextDataMatches := nextDataR.FindSubmatch([]byte(buf.String()))

  if len(nextDataMatches) != 2 {
    return NikeScrapedData{}, fmt.Errorf("'url' is invalid")
  }

  reader := bytes.NewReader(nextDataMatches[1])
  decoder := json.NewDecoder(reader)

  var nextData NextDataData
  if err := decoder.Decode(&nextData); err != nil {
    // it prob means link is invalid or some
    return NikeScrapedData{}, err
  }

  for k := range nextData.Props.PageProps.InitialState.Threads.Products {
    product, ok := nextData.Props.PageProps.InitialState.Threads.Products[k]
    if !ok {
      panic("could not get value from a map")
    }

    return NikeScrapedData {
      Title: product.Title,
      ID: nextData.Query.PBID,
      Price: product.CurrentPrice,
    }, nil
  }

  return NikeScrapedData{}, fmt.Errorf("'url' is invalid")
}
