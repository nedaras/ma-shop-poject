package handlers

import (
	"errors"
	"fmt"
	"nedas/shop/src/components"
	"nedas/shop/src/views"
	"strings"

	"github.com/labstack/echo/v4"
)

// use error codes broo
func HandleSearch(c echo.Context) error {
  url, err := convertLink(c.FormValue("url"))
  if err != nil {
    switch {
    case errors.Is(err, ErrInvalidURL):
      return render(c, components.Search("link is invalid"))
    default:
      return err
    }
  }

  d, err := scrapeNikeURL(url);
  if err != nil {
    switch {
    case errors.Is(err, ErrInvalidURL):
      return render(c, components.Search("link is invalid"))
    default:
      return render(c, components.Search("AAAAAAAAAAA"))
    }
  }

  sc, err := getSneakerContext(d, true)
  if err != nil {
    switch {
    case errors.Is(err, ErrInvalidURL):
      return render(c, components.Search("link is invalid"))
    default:
      return render(c, components.Search("AAAAAAAAAAAAAAAAAAAAAA"))
    }
  }

  c.Response().Header().Add("HX-Push-Url", fmt.Sprintf("/%s/%s", d.Slug, d.ID))
  return render(c, views.Sneaker(sc)) // todo: make sneaker a component broski
}

// need to unit test
func convertLink(str string) (string, error) {
  if str == "" {
    return "", ErrInvalidURL
  }

  i := 0;
  var b strings.Builder

  switch str[0] {
    case 'h': {
      {
        flag := "http"
        if i + len(flag) > len(str) {
          return "",ErrInvalidURL  
        }

        if str[i:i+len(flag)] != flag {
          return "",ErrInvalidURL  
        }
        i += len(flag)

        if i + 1 > len(str) {
          return "",ErrInvalidURL  
        }

        if str[i] == 's' {
          i++
        }
      }
      {
        flag := "://"
        if i + len(flag) > len(str) {
          return "", ErrInvalidURL  
        }

        if str[i:i+len(flag)] != flag {
          return "", ErrInvalidURL  
        }
        i += len(flag)
      }
    }
    fallthrough
    case 'w': {
      flag := "www."
      if i + len(flag) > len(str) {
          return "", ErrInvalidURL 
      }

      if str[i:i+len(flag)] == flag {
        i += len(flag)
      } else if i == 0 {
          return "", ErrInvalidURL 
      }
    }
    fallthrough
    case 'n': {
      flag := "nike.com/"
      if i + len(flag) > len(str) {
          return "", ErrInvalidURL 
      }

      if str[i:i+len(flag)] != flag {
        return "", ErrInvalidURL
      }
      i += len(flag)
    }
  }

  if _, err := b.WriteString("https://www.nike.com/"); err != nil {
    return "", err
  }

  if i + 2 > len(str) {
    return "", ErrInvalidURL
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
    return "", ErrInvalidURL
  }

  if str[i:i+2] != "u/" {
    return "", ErrInvalidURL
  }

final:
  i += 2
  if _, err := b.WriteString("gb/u/"); err != nil {
    return "", err
  }

  if _, err := b.WriteString(str[i:]); err != nil {
    return "", err
  }

  return b.String(), nil
}
