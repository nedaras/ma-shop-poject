package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"nedas/shop/src/components"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

// f nike https://api.nike.com/customization/builderaggregator/v2/builder/GB/en_GB/af1mid365ho22
type BuilderData struct {
  MarketingComponents []struct {
    Type string `json:"type"`
    Questions []struct {
      Type string `json:"type"`
      Answers []struct {
        Type string `json:"type"`
        DisplayName string `json:"displayName"`
        Questions []struct {
          Type string `json:"type"`
          Answers []struct {
            Type string `json:"type"`
            KeyValues []struct {
              Key string `json:"key"`
              Value string `json:"value"`
            } `json:"keyValues"`
          } `json:"answers"`
        }
      } `json:"answers"`
    } `json:"questions"`
  } `json:"marketingComponents"`
}

var (
  ErrInvalidGender = errors.New("gender is invalid")
  ErrSizesNotFound = errors.New("sizes not found")
)

func translateBuilderSizeData(d *BuilderData, men bool) ([]string, error) {
  ms := "Women's"
  if men {
    ms = "Men's"
  }

  for _, c := range d.MarketingComponents {
    if c.Type != "Size" {
      continue
    }
    for _, q1 := range c.Questions {
      if q1.Type != "Gender" {
        continue
      }
      for _, a1 := range q1.Answers {
        if a1.Type != "Gender" {
          continue
        }
        if a1.DisplayName != ms {
          continue
        }
        for _, q2 := range a1.Questions {
          if q2.Type != "Sz" {
            continue
          }
          sizes := make([]string, len(q2.Answers))
          i := 0

          for _, a2 := range q2.Answers {
            if a2.Type != "Size" {
              // idk if we should continue mb err
              continue
            }
            for _, kv := range a2.KeyValues {
              if kv.Key == "uk" {
                sizes[i] = kv.Value
                i++
              }
            }
          }

          if i == 0 {
            return []string{}, ErrSizesNotFound
          }

          if len(sizes) > i {
            sizes = sizes[:i]
          }
            
          return sizes, nil
        }
      }
    }
  }

  return []string{}, ErrSizesNotFound
}

// todo: chage to men bool
func GetSizes(p string, g string) ([]string, error) {
  if g != "men" && g != "women" {
    return []string{}, ErrInvalidGender
  }

  res, err := http.Get(fmt.Sprintf("https://api.nike.com/customization/builderaggregator/v2/builder/GB/en_GB/%s", p))
  if err != nil {
    return []string{}, err
  }
  defer res.Body.Close()

  data := &BuilderData{}
  decoder := json.NewDecoder(res.Body) 

  if err := decoder.Decode(data); err != nil {
    return []string{}, err
  }

  return translateBuilderSizeData(data, g == "men")
}

func HandleSizes(c echo.Context) error {
  gender := strings.ToLower(c.QueryParam("gender"))
  path := c.Param("path")

  if gender == "" {
    return newHTTPError(http.StatusBadRequest, "query param 'gender' is not specified");
  }

  s, err := GetSizes(path, gender)
  if err != nil {
    switch {
    case errors.Is(err, ErrInvalidGender):
      return newHTTPError(http.StatusBadRequest, "query param 'gender' is invalid");
    default:
      return err
    }
  }
  return render(c, components.Sizes(s, gender == "men"))
}
