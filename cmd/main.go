package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

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
	e.GET("/address", handlers.HandleAddress)
	e.GET("/:thread_id/:mid", handlers.HandleSneaker)

	e.POST("/htmx/search", handlers.HandleSearch)
	e.POST("/htmx/address/validate", handlers.HandleAddressValidate)
	e.POST("/htmx/product/decrement", handlers.HandleDecrement)
	e.POST("/htmx/product/increment", handlers.HandleIncrement)
	e.GET("/htmx/sizes/:path", handlers.HandleSizes)

	e.PUT("/htmx/product", handlers.HandleProduct)
	e.DELETE("/htmx/product", handlers.HandleProduct)

	// change echo error handler would better error pages

	e.Logger.Fatal(e.Start(":3000"))
}
