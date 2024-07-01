package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"nedas/shop/src/components"
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
			ViewCode               string `json:"viewCode"`
			ViewNumber             string `json:"viewNumber"`
			ImageSourceURL         string `json:"imageSourceURL"`
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
		Mid  string `json:"mid"`
		Slug string `json:"slug"`
	} `json:"query"`
}

type ProductData struct {
	Title        string  `json:"title"`
	CurrentPrice float64 `json:"currentPrice"`
	PathName     string  `json:"pathName"`
}

type NikeScrapedData struct {
	ID       string
	Slug     string
	Title    string
	Price    float64
	PathName string
}

var (
	ErrImageNotFound = errors.New("image not found")
	ErrInvalidURL    = errors.New("url is invalid")
)

// https://api.nike.com/cic/grand/v1/graphql/getfulfillmenttypesofferings/v4?variables=%7B%22countryCode%22%3A%22GB%22%2C%22currency%22%3A%22GBP%22%2C%22locale%22%3A%22en-GB%22%2C%22locationId%22%3A%22%22%2C%22locationType%22%3A%22STORE_VIEWS%22%2C%22offeringTypes%22%3A%5B%22SHIP%22%5D%2C%22postalCode%22%3A%22%22%2C%22productId%22%3A%2210c70f8d-07e3-5653-b02c-bae0e5671a45%22%7D
// some idea is we can pass threadId into a path this way the page will load fast, who cares about the search slow
func HandleSneaker(c echo.Context) error {
	path := c.Param("path")
	id := c.Param("id")

	url := fmt.Sprintf("https://www.nike.com/gb/u/%s?mid=%s", path, id)

	d, err := scrapeNikeURL(url)
	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidURL):
			return render(c, views.SimpleError(http.StatusNotFound))
		default:
			return render(c, views.SimpleError(http.StatusInternalServerError))
		}
	}

	sc, err := getSneakerContext(d, true)
	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidURL):
			return render(c, views.SimpleError(http.StatusNotFound)) // add rend err funtion
		default:
			return render(c, views.SimpleError(http.StatusInternalServerError))
		}
	}

	return render(c, views.Sneaker(sc))
}

func getSneakerContext(d NikeScrapedData, men bool) (components.SneakerContext, error) {
	ch := make(chan ErrResult, 2)

	go func() {
		val, err := getNikeConsumerData(d.ID)
		ch <- ErrResult{
			Val: val,
			Err: err,
		}
	}()

	go func() {
		val, err := GetSizes(d.PathName, men)
		ch <- ErrResult{
			Val: val,
			Err: err,
		}
	}()

	var cd *NikeConsumerData
	var s []string

	for range 2 {
		res := <-ch
		if res.Err != nil {
			return components.SneakerContext{}, res.Err
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
		return components.SneakerContext{}, errors.Join(fmt.Errorf("could net get image with id 'B'"), err)
	}

	sc := components.SneakerContext{
		Title:    d.Title,
		Price:    d.Price,
		ImageSrc: img,
		Sizes:    s,
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
	return "", errors.Join(ErrInvalidURL, ErrImageNotFound)
}

func getNikeConsumerData(id string) (*NikeConsumerData, error) {
	res, err := http.Get(fmt.Sprintf("https://api.nike.com/customization/consumer_designs/v1?filter=shortId(%s)", id))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	data := &NikeConsumerData{}
	decoder := json.NewDecoder(res.Body)

	if err := decoder.Decode(&data); err != nil {
		return nil, err
	}
	return data, nil
}

// make this take in url
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
		return NikeScrapedData{}, ErrInvalidURL
	}

	reader := bytes.NewReader(nextDataMatches[1])
	decoder := json.NewDecoder(reader)

	var nextData NextDataData
	if err := decoder.Decode(&nextData); err != nil {
		return NikeScrapedData{}, err
	}

	for k := range nextData.Props.PageProps.InitialState.Threads.Products {
		product, ok := nextData.Props.PageProps.InitialState.Threads.Products[k]
		if !ok {
			panic("could not get value from a map")
		}

		return NikeScrapedData{
			Title:    product.Title,
			ID:       nextData.Query.Mid,
			Slug:     nextData.Query.Slug,
			Price:    product.CurrentPrice,
			PathName: product.PathName,
		}, nil
	}

	return NikeScrapedData{}, ErrInvalidURL
}
