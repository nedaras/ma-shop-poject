package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"nedas/shop/pkg/handlers"
	"nedas/shop/pkg/middlewares"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	e := echo.New()

	e.Static("/", "public")

	e.Use(middleware.Logger())
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
	e.GET("/htmx/sizes/:path", handlers.HandleSizes)

	// change echo error handler would better error pages

	e.Logger.Fatal(e.Start(":3000"))
}
