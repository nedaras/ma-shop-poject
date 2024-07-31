package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"nedas/shop/pkg/apis"
	"nedas/shop/pkg/handlers"
	"nedas/shop/pkg/middlewares"
	"nedas/shop/pkg/storage"
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

	apis.SetAddressValidator(&apis.Here{})

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

	e.POST("/htmx/search", handlers.HandleSearch)
	e.POST("/htmx/add_to_bag", handlers.AddToBag)
	e.POST("/htmx/checkout", handlers.HandleCheckout)
	e.POST("/htmx/login", handlers.HandleLogin)
	e.POST("/htmx/logout", handlers.HandleLogout)
	e.POST("/htmx/product/decrement", handlers.HandleDecrement)
	e.POST("/htmx/product/increment", handlers.HandleIncrement)
	e.GET("/htmx/address", handlers.HandleAddressEditor)
	e.GET("/htmx/address/:id", handlers.HandleAddressEditor)
	e.GET("/htmx/sizes/:path", handlers.HandleSizes)

	e.PUT("/htmx/product", handlers.HandleProduct)
	e.PUT("/htmx/address/:id", handlers.HandlePutAddress)

	e.DELETE("/htmx/product", handlers.HandleProduct)
	e.DELETE("/htmx/address/:id", handlers.HandleDeleteAddress)

	// todo: we need like meta tags for some nive embeds
	// todo: we prob will drop cassandra it aint a db for cuz we kinda will be doing joins manually
	// todo: change echo error handler would better error pages

	e.Logger.Fatal(e.Start(":3000"))
}
