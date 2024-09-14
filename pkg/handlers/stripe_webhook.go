package handlers

import (
	"encoding/json"
	"io"
	"nedas/shop/pkg/utils"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/stripe/stripe-go/v79"
	"github.com/stripe/stripe-go/v79/checkout/session"
	"github.com/stripe/stripe-go/v79/webhook"
)

// users shoud have list of order_ids
// db: order_id -> ...product, ...address, status
func HandleStripeWebhook(c echo.Context) error {
	storage := getStorage(c)

	payload, err := io.ReadAll(c.Request().Body) // todo: add max bytes
	if err != nil {
		utils.Logger().Error(err)
		return err
	}

	event, err := webhook.ConstructEvent(payload, c.Request().Header.Get("Stripe-Signature"), utils.Getenv("STRIPE_SIGNING_SECRET"))
	if err != nil {
		utils.Logger().Error(err)
		return err
	}

	if event.Type == "checkout.session.completed" {
		var session stripe.CheckoutSession
		err := json.Unmarshal(event.Data.Raw, &session)
		if err != nil {
			utils.Logger().Error(err)
			return err
		}

		userId := session.Metadata["user_id"]
		items := getLineItems(session.ID)

		// we need async loop or smth
		for _, item := range items.Data {
			tid := item.Price.Metadata["tid"]
			mid := item.Price.Metadata["mid"]
			size := item.Price.Metadata["size"]

			product, err := getProduct(tid + ":" + mid)
			if err != nil {
				return err
			}

			// todo: we should check if database has other products so we would not remove like all products
			if err := storage.DeleteProduct(userId, tid, mid, size); err != nil {
				_ = product
				utils.Logger().Error(err)
				return err
			}
		}
	}

	return c.NoContent(http.StatusOK)
}

func getLineItems(id string) *stripe.LineItemList {
	params := &stripe.CheckoutSessionListLineItemsParams{
		Session: stripe.String(id),
	}

	return session.ListLineItems(params).LineItemList()
}
