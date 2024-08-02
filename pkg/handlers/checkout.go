package handlers

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func HandleCheckout(c echo.Context) error {
	session := getSession(c)
	storage := getStorage(c)

	if session == nil {
		return unauthorized(c)
	}

	// todo: note we have to make like a cache for what we tryna checkout cuz we dont want any out of sync problems
	products, err := getProducts(session.UserId, storage)
	if err != nil {
		return err
	}

	if len(products) == 0 {
		return redirect(c, "/bag")
	}

	for _, p := range products {
		// todo: we wanna make this precise like idk some math module and convert the amount float to string
		fmt.Println("price: ", float64(p.Amount)*p.Product.Price)
	}

	return c.NoContent(http.StatusNotFound)
}
