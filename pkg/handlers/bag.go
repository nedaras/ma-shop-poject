package handlers

import (
	"encoding/json"
	"errors"
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
	products            = []string{"bb049e5fc-e1a4-4196-92c3-439ed3c475d1:3475937855", "e3864a31-60d8-470a-8f62-41cc7c0688bd:4063348121"}
	ErrInvalidProductID = errors.New("invalid product id")
  ErrNotFound = errors.New("could not found requested resource")
)

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
		return views.Product{}, ErrInvalidProductID
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
				Title: p.ProductContent.Title,
				Price: p.MerchPrice.CurrentPrice,
				Image: img,
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
			ch <- struct{
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

  // a b c _ f _ _ _ g g e _
  // _ _ _ a _ b _ c d g
  // _ a
  // can i like in O(n) like remove white space from products, without making a new array
  // we fr fr need to unit test this one

  fe := -1
  for i, v := range products {
    if fe == size {
      break
    }

    if v == (views.Product{}) {
      if fe == -1 {
        fe = i
      }
      continue
    }

    if fe == -1 {
      continue
    }

    products[fe] = v
    products[i] = views.Product{}
    fe++
  }

  return products[0:size], nil
}
