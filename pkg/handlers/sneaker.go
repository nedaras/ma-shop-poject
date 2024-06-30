package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"nedas/shop/src/views"
	"net/http"
	"regexp"
	"strings"

	"github.com/labstack/echo/v4"
)

// use the template to load smaller images
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

// how to err if like query is undefined hu?
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
    //PBID string `json:"pbid"`
    Mid string `json:"mid"`
    Slug string `json:"slug"`
  } `json:"query"`
}

type ProductData struct {
  Title string `json:"title"`
  CurrentPrice float64 `json:"currentPrice"`
  PathName string `json:"pathName"`
}

type NikeScrapedData struct {
  ID string
  Slug string
  Title string
  Price float64
  PathName string
}

var (
  ErrImageNotFound = errors.New("image not found")
)

func HandleSneaker(c echo.Context) error {
  path := c.Param("path")
  id := c.Param("id")

  d, err := scrapeNikeURL(path, id);
  if err != nil {
    return err // ret 404
  }

  sc, err := getSneakerContext(d, true)
  if err != nil {
    return err // ret hard 500
  }
  return render(c, views.Sneaker(sc))
}

// and a way to edit nike image size

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

func getSneakerContext(d NikeScrapedData, men bool) (views.SneakerContext, error) {
  ch := make(chan ErrResult, 2)

  go func(id string, ch chan<- ErrResult) {
    val, err := getNikeConsumerData(id)
    ch <- ErrResult {
      Val: val,
      Err: err,
    }
  }(d.ID, ch)

  go func(p string, men bool, ch chan<- ErrResult) {
    val, err := GetSizes(p, men)
    ch <- ErrResult {
      Val: val,
      Err: err,
    }
  }(d.PathName, men, ch)

  var cd *NikeConsumerData
  var s []string

  for range(2) {
    res := <- ch
    if res.Err != nil {
      return views.SneakerContext{}, res.Err
    }

    switch v := res.Val.(type) {
    case *NikeConsumerData:
      cd = v
    case []string:
      s = v
    default:
      panic("got invalid type")
    }
  }

  img, err := getImageByID(cd, "B")
  if err != nil {
    return views.SneakerContext{}, errors.Join(fmt.Errorf("could net get image with id 'B'"), err)
  }

  sc := views.SneakerContext{
    Title: d.Title,
    Price: d.Price,
    ImageSrc: img,
    Sizes: s,
    PathName: d.PathName,
  }

  return sc, nil
}

func getImageByID(d *NikeConsumerData, id string) (string, error) {
  for _, i := range d.Objects {
    for _, image := range i.Imagery {
      if image.ViewCode == id {
        return image.ImageSourceURL, nil
      }
    }
  }
  return "", ErrImageNotFound
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

func scrapeNikeURL(path string, id string) (NikeScrapedData, error) {
  nextDataR := regexp.MustCompile(`<script id="__NEXT_DATA__" type="application\/json">(.+)<\/script>`)
  url := fmt.Sprintf("https://www.nike.com/gb/u/%s?mid=%s", path, id)

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
      ID: nextData.Query.Mid,
      Slug: nextData.Query.Slug,
      Price: product.CurrentPrice,
      PathName: product.PathName,
    }, nil
  }

  return NikeScrapedData{}, fmt.Errorf("'url' is invalid")
}
