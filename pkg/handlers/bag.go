package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"nedas/shop/src/views"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

type ProductFeedData struct {
	Objects []struct {
		ProductInfo []struct {
			MerchPrice struct {
				CurrentPrice float64 `json:"currentPrice"`
			} `json:"merchPrice"`
			ProductContent struct {
				Title string `json:"title"`
			} `json:"productContent"`
			CustomizedPreBuild struct {
				Legacy struct {
					PathName string `json:"pathName"`
				} `json:"legacy"`
			} `json:"customizedPreBuild"`
		} `json:"productInfo"`
	} `json:"objects"`
}

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

// placeholder, what db we will use
var (
	products    = []string{"b049e5fc-e1a4-4196-92c3-439ed3c475d1:3475937855", "e3864a31-60d8-470a-8f62-41cc7c0688bd:4063348121"}
	ErrNotFound = errors.New("could not found requested resource")
)

func HandleBag(c echo.Context) error {
	products, err := getProducts(products)
	if err != nil {
		return err
	}

	bc := views.BagContext{
		Products: products,
	}

	return render(c, views.Bag(bc))
}

// Any returned error will be of type [*NikeAPIError].
func getProductFeedData(tid string) (*ProductFeedData, error) {
	url := "https://api.nike.com/product_feed/rollup_threads/v2?filter=marketplace(GB)&filter=language(en-GB)&filter=employeePrice(true)&filter=id(" + tid + ")&consumerChannelId=d9a5bc42-4b9c-4976-858a-f159cf99c647"
	res, err := http.Get(url)

	if err != nil {
		return nil, &NikeAPIError{URL: url, Err: err}
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		switch res.StatusCode {
		case 404:
			return nil, &NikeAPIError{URL: url, Err: ErrNotFound}
		default:
			return nil, &NikeAPIError{URL: url, Err: fmt.Errorf("got unexpected response code '%d'", res.StatusCode)}
		}
	}

	if res.Header.Get("Content-Type") != "application/json" {
		return nil, &NikeAPIError{URL: url, Err: errors.New("responded content is not in json form")}
	}

	data := &ProductFeedData{}
	decoder := json.NewDecoder(res.Body)

	if err := decoder.Decode(&data); err != nil {
		return nil, &NikeAPIError{URL: url, Err: err}
	}

	if len(data.Objects) == 0 {
		return nil, &NikeAPIError{URL: url, Err: ErrNotFound}
	}

	return data, nil
}

// Any returned error will be of type [*NikeAPIError].
func getProduct(id string) (views.Product, error) {
	arr := strings.SplitN(id, ":", 2)
	if len(arr) != 2 {
		panic("passed string is not split by ':'")
	}

	tid, mid := arr[0], arr[1]
	ch := make(chan ErrResult[any], 2)

	var (
		pf *ProductFeedData
		cd *NikeConsumerData
	)

	go func() {
		val, err := getProductFeedData(tid)
		ch <- ErrResult[any]{
			Val: val,
			Err: err,
		}
	}()

	go func() {
		val, err := getNikeConsumerData(mid)
		ch <- ErrResult[any]{
			Val: val,
			Err: err,
		}
	}()

	for range 2 {
		res := <-ch
		if res.Err != nil {
			return views.Product{}, res.Err
		}
		switch v := res.Val.(type) {
		case *NikeConsumerData:
			cd = v
		case *ProductFeedData:
			pf = v
		default:
			panic("got invalid type")
		}
	}

	img := getImageByID(cd, "B")
	if img == "" {
		return views.Product{}, &NikeAPIError{
			URL: "https://api.nike.com/customization/consumer_designs/v1?filter=shortId(" + mid + ")",
			Err: ErrNotFound,
		}
	}

	for _, o := range pf.Objects {
		for _, p := range o.ProductInfo {
			if p.ProductContent.Title == "" {
				continue
			}
			if p.MerchPrice.CurrentPrice == 0.0 {
				continue
			}
			return views.Product{
				Title:    p.ProductContent.Title,
				Price:    p.MerchPrice.CurrentPrice,
				Image:    img,
				PathName: p.CustomizedPreBuild.Legacy.PathName,
			}, nil
		}
	}

	return views.Product{}, &NikeAPIError{
		URL: "https://api.nike.com/product_feed/rollup_threads/v2?filter=marketplace(GB)&filter=language(en-GB)&filter=employeePrice(true)&filter=id(" + tid + ")&consumerChannelId=d9a5bc42-4b9c-4976-858a-f159cf99c647",
		Err: ErrNotFound,
	}
}

// Any returned error will be of type [*NikeAPIError].
func getProducts(p []string) ([]views.Product, error) {
	if len(p) == 1 {
		p, err := getProduct(p[0])
		if err != nil {
			return []views.Product{}, err
		}
		return []views.Product{p}, nil
	}

	ch := make(chan struct {
		i int
		p views.Product
		e error
	}, len(p))

	products := make([]views.Product, len(p))
	size := 0

	for i, id := range p {
		go func() {
			val, err := getProduct(id)
			ch <- struct {
				i int
				p views.Product
				e error
			}{
				i: i,
				p: val,
				e: err,
			}
		}()
	}

	for range p {
		res := <-ch
		if res.e != nil {
			if errors.Is(res.e, ErrNotFound) {
				continue
			}
			return []views.Product{}, res.e
		}
		products[res.i] = res.p
		size++
	}

	if size == len(p) {
		return products, nil
	}

	strip(&products)
	return products, nil
}

func strip[T comparable](arr *[]T) {
	var mt T
	fe := -1

	for i, v := range *arr {
		if v == mt {
			if fe == -1 {
				fe = i
			}
			continue
		}

		if fe == -1 {
			continue
		}

		(*arr)[fe] = v
		(*arr)[i] = mt
		fe++
	}

	if fe != -1 {
		*arr = (*arr)[:fe]
	}
}

// Any returned error will be of type [*NikeAPIError].
func getNikeConsumerData(mid string) (*NikeConsumerData, error) {
	url := "https://api.nike.com/customization/consumer_designs/v1?filter=shortId(" + mid + ")"
	res, err := http.Get(url)
	if err != nil {
		return nil, &NikeAPIError{URL: url, Err: err}
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		switch res.StatusCode {
		case 404:
			return nil, &NikeAPIError{URL: url, Err: ErrNotFound}
		default:
			return nil, &NikeAPIError{URL: url, Err: fmt.Errorf("got unexpected response code '%d'", res.StatusCode)}
		}
	}

	if res.Header.Get("Content-Type") != "application/json" {
		return nil, &NikeAPIError{URL: url, Err: errors.New("responded content is not in json form")}
	}

	data := &NikeConsumerData{}
	decoder := json.NewDecoder(res.Body)

	if err := decoder.Decode(&data); err != nil {
		return nil, &NikeAPIError{URL: url, Err: err}
	}

	if len(data.Objects) == 0 {
		return nil, &NikeAPIError{URL: url, Err: ErrNotFound}
	}

	if len(data.Errors) != 0 {
		return nil, &NikeAPIError{URL: url, Err: fmt.Errorf("got some unexpected errors %v", data.Errors...)}
	}

	return data, nil
}

func getImageByID(d *NikeConsumerData, id string) string {
	for _, i := range d.Objects {
		for _, image := range i.Imagery {
			if image.ViewCode == id {
				return image.ImageSourceURL
			}
		}
	}
	return ""
}

type NikeAPIError struct {
	URL string
	Err error
}

func (e *NikeAPIError) Error() string {
	return "'" + e.URL + "': " + e.Err.Error()
}

func (e *NikeAPIError) Unwrap() error {
	return e.Err
}
