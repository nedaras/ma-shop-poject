package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/stripe/stripe-go/v79"

	"nedas/shop/pkg/apis"
	"nedas/shop/pkg/handlers"
	"nedas/shop/pkg/middlewares"
	"nedas/shop/pkg/storage"
	"nedas/shop/pkg/utils"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	storage, err := storage.NewCassandra()
	if err != nil {
		log.Fatal(err)
	}
	defer storage.Close()

	apis.SetAddressValidator(apis.NewHere(1000))

	stripe.Key = utils.Getenv("STRIPE_SECRET_KEY")

	e := echo.New()
	e.Static("/", "public")

	e.Use(middleware.Logger())
	e.Use(middlewares.StorageMiddleware(storage))
	e.Use(middlewares.AuthMiddleware)

	e.GET("/", handlers.HandleIndex)
	e.GET("/login", handlers.HandleLogin)
	e.GET("/login/google", handlers.HandleGoogleLogin)
	e.GET("/bag", handlers.HandleBag)
	e.GET("/account", handlers.HandleAccount)
	e.GET("/addresses", handlers.HandleAddresses)
	e.GET("/:thread_id/:mid", handlers.HandleSneaker)
	e.GET("/address/:id", handlers.HandleAddress)

	e.POST("/htmx/search", handlers.HandleSearch)
	e.POST("/htmx/add_to_bag", handlers.AddToBag)
	e.POST("/htmx/login", handlers.HandleLogin)
	e.POST("/htmx/logout", handlers.HandleLogout)
	e.POST("/htmx/checkout", handlers.HandleCheckout)
	e.POST("/htmx/product/decrement", handlers.HandleDecrement)
	e.POST("/htmx/product/increment", handlers.HandleIncrement)
	e.GET("/htmx/sizes/:path", handlers.HandleSizes)
	e.GET("/htmx/address", handlers.HandleCreateAddress)

	e.PUT("/htmx/product", handlers.HandleProduct)
	e.PUT("/htmx/address/:id", handlers.HandlePutAddress)

	e.DELETE("/htmx/product", handlers.HandleProduct)
	e.DELETE("/htmx/address/:id", handlers.HandleDeleteAddress)

	e.HTTPErrorHandler = handlers.ErrorHandler
	e.Logger.Fatal(e.Start(":3000"))
}
