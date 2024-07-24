package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"nedas/shop/pkg/models"
	"nedas/shop/pkg/storage"
	"nedas/shop/src/components"
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
				//Subtitle string `json:"subtitle"`
			} `json:"productContent"`
			CustomizedPreBuild struct {
				Groups []struct {
					Legacy struct {
						PIID string `json:"piid"`
						Slug string `json:"slug"`
					} `json:"legacy"`
				} `json:"groups"`
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

var (
	ErrNotFound = errors.New("could not found requested resource")
)

func HandleBag(c echo.Context) error {
	session := getSession(c)
	storage := getStorage(c)

	if session == nil {
		// products from cookies or sum
		return render(c, views.Bag([]components.BagProductContext{}))
	}

	// this functio one day will just take in c and do its own shit
	// err can be from nike api or from storage
	products, err := getProducts(session.UserId, storage)
	if err != nil {
		c.Logger().Error(err)
		return err
	}

	return render(c, views.Bag(products))
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
func getProduct(id string) (models.Product, error) {
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
			return models.Product{}, res.Err
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
		return models.Product{}, &NikeAPIError{
			URL: "https://api.nike.com/customization/consumer_designs/v1?filter=shortId(" + mid + ")",
			Err: ErrNotFound,
		}
	}

	for _, o := range pf.Objects {
		for _, p := range o.ProductInfo {
			if p.CustomizedPreBuild.Legacy.PathName == "" {
				continue
			}
			for _, g := range p.CustomizedPreBuild.Groups {
				if g.Legacy.Slug == "" {
					continue
				}

				slug := g.Legacy.Slug
				if g.Legacy.PIID != "" {
					slug += "-" + g.Legacy.PIID
				}

				return models.Product{
					Title:    p.ProductContent.Title,
					Price:    p.MerchPrice.CurrentPrice,
					Image:    img,
					PathName: p.CustomizedPreBuild.Legacy.PathName,
					Mid:      mid,
					ThreadId: tid,
					Slug:     slug,
				}, nil
			}
		}
	}

	return models.Product{}, &NikeAPIError{
		URL: "https://api.nike.com/product_feed/rollup_threads/v2?filter=marketplace(GB)&filter=language(en-GB)&filter=employeePrice(true)&filter=id(" + tid + ")&consumerChannelId=d9a5bc42-4b9c-4976-858a-f159cf99c647",
		Err: ErrNotFound,
	}
}

func getProducts(userId string, storage storage.Storage) ([]components.BagProductContext, error) {
	storageProducts, err := storage.GetProducts(userId)
	if err != nil {
		return []components.BagProductContext{}, err
	}

	if len(storageProducts) == 1 {
		product := storageProducts[0]
		// todo: fr add validate product
		p, err := getProduct(product.ProductId)

		if err != nil {
			return []components.BagProductContext{}, err
		}

		amount, err := storage.GetProductAmount(userId, p.ThreadId, p.Mid, product.Size)
		if err != nil {
			return []components.BagProductContext{}, err
		}

		return []components.BagProductContext{{
			Product: p,
			Size:    product.Size,
			Amount:  amount,
		}}, nil
	}

	ch := make(chan struct {
		i       int
		product models.Product
		size    string
		amount  uint8
		err     error
	}, len(storageProducts))

	products := make([]components.BagProductContext, len(storageProducts))
	size := 0

	for i, product := range storageProducts {
		go func() {
			// todo: frr add validate
			val, err := getProduct(product.ProductId)
			if err != nil {
				ch <- struct {
					i       int
					product models.Product
					size    string
					amount  uint8
					err     error
				}{
					i:       i,
					product: models.Product{},
					size:    product.Size,
					amount:  0,
					err:     err,
				}
				return
			}

			amount, err := storage.GetProductAmount(userId, val.ThreadId, val.Mid, product.Size)
			ch <- struct {
				i       int
				product models.Product
				size    string
				amount  uint8
				err     error
			}{
				i:       i,
				product: val,
				size:    product.Size,
				amount:  amount,
				err:     err,
			}

		}()
	}

	for range storageProducts {
		res := <-ch
		if res.err != nil {
			if errors.Is(res.err, ErrNotFound) {
				continue
			}
			return []components.BagProductContext{}, res.err
		}
		products[res.i] = components.BagProductContext{
			Product: res.product,
			Size:    res.size,
			Amount:  res.amount,
		}
		size++
	}

	if size == len(storageProducts) {
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
