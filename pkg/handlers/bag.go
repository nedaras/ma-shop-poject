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
		} `json:"productInfo"`
	} `json:"objects"`
}

// placeholder, what db we will use
var (
	products            = []string{"b049e5fc-e1a4-4196-92c3-439ed3c475d1:3475937855", "e3864a31-60d8-470a-8f62-41cc7c0688bd:4063348121"}
	ErrInvalidProductID = errors.New("invalid product id")
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

func getProductFeedData(tid string) (*ProductFeedData, error) {
	url := fmt.Sprintf("https://api.nike.com/product_feed/rollup_threads/v2?filter=marketplace(GB)&filter=language(en-GB)&filter=employeePrice(true)&filter=id(%s)&consumerChannelId=d9a5bc42-4b9c-4976-858a-f159cf99c647", tid)

	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	data := &ProductFeedData{}
	decoder := json.NewDecoder(res.Body)

	if err := decoder.Decode(&data); err != nil {
		return nil, err
	}

	if len(data.Objects) == 0 {
		return nil, ErrInvalidProductID
	}

	return data, nil
}

func getProduct(id string) (views.Product, error) {
	arr := strings.SplitN(id, ":", 2)
	if len(arr) != 2 {
		return views.Product{}, ErrInvalidProductID
	}

	tid := arr[0]
	mid := arr[1]

	ch := make(chan ErrResult[any], 2)

	var cd *NikeConsumerData
	var pf *ProductFeedData

	go func() {
		val, err := getProductFeedData(tid)
		if err != nil {
			err = errors.Join(fmt.Errorf("could net get product feed data with thread id '%s'", tid), err)
		}
		ch <- ErrResult[any]{
			Val: val,
			Err: err,
		}
	}()

	go func() {
		val, err := getNikeConsumerData(mid)
		if err != nil {
			err = errors.Join(fmt.Errorf("could not get consumer data with mid '%s'", mid), err)
		}
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

	img, err := getImageByID(cd, "B")
	if err != nil {
		return views.Product{}, errors.Join(fmt.Errorf("could net get image with id 'B'"), ErrInvalidProductID, err)
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
				Title: p.ProductContent.Title,
				Price: p.MerchPrice.CurrentPrice,
				Image: img,
			}, nil
		}
	}
	return views.Product{}, ErrInvalidProductID
}

func getProducts(p []string) ([]views.Product, error) {
	if len(p) == 1 {
		p, err := getProduct(p[0])
		if err != nil {
			return []views.Product{}, err
		}
		return []views.Product{p}, nil
	}

	type indexed = struct {
		i int
		p views.Product
		e error
	}

	ch := make(chan indexed, len(products))
	products := make([]views.Product, len(products))

	for i, id := range p {
		go func() {
			val, err := getProduct(id)
			ch <- indexed{
				i: i,
				p: val,
				e: err,
			}
		}()
	}

	for range p {
		res := <-ch
		if res.e != nil {
			// idk what todo cuz if one errors what er all err then, we need to think
			return []views.Product{}, res.e
		}
		products[res.i] = res.p
	}
	return products, nil
}
