package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"nedas/shop/pkg/models"
	"nedas/shop/src/components"
	"nedas/shop/src/views"
	"net/http"
	"regexp"
	"strings"

	"github.com/labstack/echo/v4"
)

type NextData struct {
	Props struct {
		PageProps struct {
			InitialState struct {
				Threads struct {
					Products map[string]struct {
						Title        string  `json:"title"`
						CurrentPrice float64 `json:"currentPrice"`
						PathName     string  `json:"pathName"`
						ThreadId     string  `json:"threadId"`
					} `json:"products"`
				} `json:"Threads"`
			} `json:"initialState"`
		} `json:"pageProps"`
	} `json:"props"`
	Query struct {
		Mid  string `json:"mid"` // mb use metricId from "__NEXT_DATA__.props.pageProps.initialState.NikeId"
		Slug string `json:"slug"`
	} `json:"query"`
}

type NikeScrapedData struct {
	Mid      string
	Title    string
	Price    float64
	PathName string
	ThreadId string
	Slug     string
}

var (
	NextDataRegexp = regexp.MustCompile(`<script id="__NEXT_DATA__" type="application\/json">(.+)<\/script>`)
)

func HandleSearch(c echo.Context) error {
	url := convertLink(c.FormValue("url"))
	if url == "" {
		return renderWithStatus(http.StatusNotFound, c, components.Search("Nike By You link is invalid."))
	}

	data, err := scrapeNikeURL(url)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return renderWithStatus(http.StatusNotFound, c, components.Search("Nike By You link is invalid."))
		}
		c.Logger().Error(err)
		return renderWithStatus(http.StatusInternalServerError, c, components.Search("Something went wrong. Please try again later."))
	}

	ch := make(chan ErrResult[any], 2)

	var (
		img   string
		sizes []string
	)

	go func() {
		val, err := getImageWithId(data.Mid, "B")
		ch <- ErrResult[any]{
			Val: val,
			Err: err,
		}
	}()

	go func() {
		val, err := GetSizes(data.PathName, true)
		ch <- ErrResult[any]{
			Val: val,
			Err: err,
		}
	}()

	for range 2 {
		res := <-ch
		if res.Err != nil {
			if errors.Is(res.Err, ErrNotFound) {
				return renderWithStatus(http.StatusNotFound, c, components.Search("Nike By You link is invalid."))
			}
			c.Logger().Error(res.Err)
			return renderWithStatus(http.StatusInternalServerError, c, components.Search("Something went wrong. Please try again later."))
		}
		switch v := res.Val.(type) {
		case string:
			img = v
		case []string:
			sizes = v
		default:
			panic("got invalid type")
		}
	}

	product := models.Product{
		Title:    data.Title,
		Price:    data.Price,
		PathName: data.PathName,
		Mid:      data.Mid,
		ThreadId: data.ThreadId,
		Slug:     data.Slug,
		Image:    img,
	}

	c.Response().Header().Add("HX-Push-Url", "/"+data.ThreadId+"/"+data.Mid)
	return render(c, views.Sneaker(views.SneakerContext{
		Product:  product,
		Sizes:    sizes,
		LoggedIn: getSession(c) != nil,
	}))
}

// Any returned error will be of type [*NikeAPIError].
func getImageWithId(mid string, id string) (string, error) {
	val, err := getNikeConsumerData(mid)
	if err != nil {
		return "", err
	}

	img := getImageByID(val, id)
	if img == "" {
		return "", &NikeAPIError{URL: "https://api.nike.com/customization/consumer_designs/v1?filter=shortId(" + mid + ")", Err: ErrNotFound}
	}

	return img, nil
}

// Any returned error will be of type [*NikeAPIError].
func scrapeNikeURL(url string) (NikeScrapedData, error) {
	res, err := http.Get(url)
	if err != nil {
		return NikeScrapedData{}, &NikeAPIError{URL: url, Err: err}
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		switch res.StatusCode {
		case 404:
			return NikeScrapedData{}, &NikeAPIError{URL: url, Err: ErrNotFound}
		default:
			return NikeScrapedData{}, &NikeAPIError{URL: url, Err: fmt.Errorf("got unexpected response code '%d'", res.StatusCode)}
		}
	}

	if res.Header.Get("Content-Type") != "text/html; charset=UTF-8" {
		return NikeScrapedData{}, &NikeAPIError{URL: url, Err: errors.New("responded content is not in UTF-8 html form")}
	}

	buf := new(strings.Builder)
	if _, err = io.Copy(buf, res.Body); err != nil {
		return NikeScrapedData{}, &NikeAPIError{URL: url, Err: err}
	}

	nextDataMatches := NextDataRegexp.FindSubmatch([]byte(buf.String()))

	if len(nextDataMatches) != 2 {
		return NikeScrapedData{}, &NikeAPIError{URL: url, Err: ErrNotFound}
	}

	reader := bytes.NewReader(nextDataMatches[1])
	decoder := json.NewDecoder(reader)

	var nextData NextData
	if err := decoder.Decode(&nextData); err != nil {
		return NikeScrapedData{}, &NikeAPIError{URL: url, Err: err}
	}

	if nextData.Query.Mid == "" {
		return NikeScrapedData{}, &NikeAPIError{URL: url, Err: ErrNotFound}
	}

	for k := range nextData.Props.PageProps.InitialState.Threads.Products {
		product := nextData.Props.PageProps.InitialState.Threads.Products[k]
		if product.PathName == "" {
			continue
		}

		return NikeScrapedData{
			Title:    product.Title,
			Mid:      nextData.Query.Mid,
			ThreadId: product.ThreadId,
			Price:    product.CurrentPrice,
			PathName: product.PathName,
			Slug:     nextData.Query.Slug,
		}, nil
	}

	return NikeScrapedData{}, &NikeAPIError{URL: url, Err: ErrNotFound}
}

// need to unit test
func convertLink(str string) string {
	if str == "" {
		return ""
	}

	i := 0
	var b strings.Builder

	switch str[0] {
	case 'h':
		{
			{
				flag := "http"
				if i+len(flag) > len(str) {
					return ""
				}

				if str[i:i+len(flag)] != flag {
					return ""
				}
				i += len(flag)

				if i+1 > len(str) {
					return ""
				}

				if str[i] == 's' {
					i++
				}
			}
			{
				flag := "://"
				if i+len(flag) > len(str) {
					return ""
				}

				if str[i:i+len(flag)] != flag {
					return ""
				}
				i += len(flag)
			}
		}
		fallthrough
	case 'w':
		{
			flag := "www."
			if i+len(flag) > len(str) {
				return ""
			}

			if str[i:i+len(flag)] == flag {
				i += len(flag)
			} else if i == 0 {
				return ""
			}
		}
		fallthrough
	case 'n':
		{
			flag := "nike.com/"
			if i+len(flag) > len(str) {
				return ""
			}

			if str[i:i+len(flag)] != flag {
				return ""
			}
			i += len(flag)
		}
	}

	b.WriteString("https://www.nike.com/")

	if i+2 > len(str) {
		return ""
	}

	if str[i:i+2] == "u/" {
		goto final
	}

	for i < len(str) {
		if str[i] == '/' {
			i++
			break
		}
		i++
	}

	if i+2 > len(str) {
		return ""
	}

	if str[i:i+2] != "u/" {
		return ""
	}

final:
	i += 2

	b.WriteString("gb/u/")
	b.WriteString(str[i:])

	return b.String()
}
