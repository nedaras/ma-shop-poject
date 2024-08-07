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
		Type      string `json:"type"`
		Questions []struct {
			Type    string `json:"type"`
			Answers []struct {
				Type        string `json:"type"`
				DisplayName string `json:"displayName"`
				Questions   []struct {
					Type    string `json:"type"`
					Answers []struct {
						Type      string `json:"type"`
						KeyValues []struct {
							Key   string `json:"key"`
							Value string `json:"value"`
						} `json:"keyValues"`
					} `json:"answers"`
				}
			} `json:"answers"`
		} `json:"questions"`
	} `json:"marketingComponents"`
}

func HandleSizes(c echo.Context) error {
	gender := strings.ToLower(c.QueryParam("gender"))
	path := c.Param("path")

	if gender == "" {
		return newHTTPError(http.StatusBadRequest, "query param 'gender' is not specified")
	}

	if gender != "men" && gender != "women" {
		return newHTTPError(http.StatusBadRequest, "query param 'gender' is invalid")
	}

	s, err := GetSizes(path, gender == "men")
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return newHTTPError(http.StatusNotFound, "could not find "+gender+"'s sizes")
		}
		return err
	}
	return render(c, components.Sizes(s, path, gender == "men"))
}

// Any returned error will be of type [*NikeAPIError].
func GetSizes(path string, men bool) ([]string, error) {
	url := "https://api.nike.com/customization/builderaggregator/v2/builder/GB/en_GB/" + path
	res, err := http.Get(url)

	if err != nil {
		return []string{}, &NikeAPIError{URL: url, Err: err}
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		switch res.StatusCode {
		case http.StatusNotFound:
			return []string{}, &NikeAPIError{URL: url, Err: ErrNotFound}
		default:
			return []string{}, &NikeAPIError{URL: url, Err: fmt.Errorf("got unexpected response code '%d'", res.StatusCode)}
		}
	}

	if res.Header.Get("Content-Type") != "application/json" {
		return []string{}, &NikeAPIError{URL: url, Err: errors.New("responded content is not in json form")}
	}

	data := &BuilderData{}
	decoder := json.NewDecoder(res.Body)

	if err := decoder.Decode(data); err != nil {
		return []string{}, &NikeAPIError{URL: url, Err: err}
	}

	gstr := "Women's"
	if men {
		gstr = "Men's"
	}

	for _, c := range data.MarketingComponents {
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
				if a1.DisplayName != gstr && a1.DisplayName != "Unisex" {
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
							continue
						}
						for _, kv := range a2.KeyValues {
							if kv.Key == "uk" {
								sizes[i] = strings.Trim(kv.Value, " \n\r\t")
								i++
							}
						}
					}

					if i == 0 {
						return []string{}, &NikeAPIError{URL: url, Err: ErrNotFound}
					}

					if len(sizes) > i {
						sizes = sizes[:i]
					}

					return sizes, nil
				}
			}
		}
	}

	return []string{}, &NikeAPIError{URL: url, Err: ErrNotFound}
}

// Any returned error will be of type [*NikeAPIError].
func GetAllSizes(path string) ([]string, error) {
	url := "https://api.nike.com/customization/builderaggregator/v2/builder/GB/en_GB/" + path
	res, err := http.Get(url)

	if err != nil {
		return []string{}, &NikeAPIError{URL: url, Err: err}
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		switch res.StatusCode {
		case 404:
			return []string{}, &NikeAPIError{URL: url, Err: ErrNotFound}
		default:
			return []string{}, &NikeAPIError{URL: url, Err: fmt.Errorf("got unexpected response code '%d'", res.StatusCode)}
		}
	}

	if res.Header.Get("Content-Type") != "application/json" {
		return []string{}, &NikeAPIError{URL: url, Err: errors.New("responded content is not in json form")}
	}

	data := &BuilderData{}
	decoder := json.NewDecoder(res.Body)

	if err := decoder.Decode(data); err != nil {
		return []string{}, &NikeAPIError{URL: url, Err: err}
	}

	for _, c := range data.MarketingComponents {
		if c.Type != "Size" {
			continue
		}
		for _, q1 := range c.Questions {
			if q1.Type != "Gender" {
				continue
			}

			length := 0
			for _, a1 := range q1.Answers {
				if a1.Type != "Gender" {
					continue
				}
				switch a1.DisplayName {
				case "Men's", "Women's", "Unisex":
					for _, q2 := range a1.Questions {
						if q2.Type != "Sz" {
							continue
						}
						length += len(q2.Answers)
					}
				}
			}
			if length == 0 {
				return []string{}, &NikeAPIError{URL: url, Err: ErrNotFound}
			}
			sizes := make([]string, length)
			for _, a1 := range q1.Answers {
				if a1.Type != "Gender" {
					continue
				}
				switch a1.DisplayName {
				case "Men's", "Women's", "Unisex":
					for _, q2 := range a1.Questions {
						if q2.Type != "Sz" {
							continue
						}
						for _, a2 := range q2.Answers {
							if a2.Type != "Size" {
								continue
							}
							for _, kv := range a2.KeyValues {
								if kv.Key == "uk" {
									sizes[len(sizes)-length] = strings.Trim(kv.Value, " \n\r\t")
									length--
								}
							}
						}
					}
				}
			}
			if length == 0 {
				return sizes, nil
			}
			return sizes[:len(sizes)-length], nil
		}
	}

	return []string{}, &NikeAPIError{URL: url, Err: ErrNotFound}
}
