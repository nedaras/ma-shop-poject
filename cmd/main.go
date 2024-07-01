package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"nedas/shop/pkg/handlers"
)

func main() {

	e := echo.New()

  e.Static("/", "public");

  e.Use(middleware.Logger())
	e.GET("/", handlers.HandleIndex)
	e.GET("/address", handlers.HandleAddress)
  e.GET("/:path/:id", handlers.HandleSneaker)

  e.POST("/htmx/search", handlers.HandleSearch)
	e.POST("/htmx/address/validate", handlers.HandleAddressValidate)
  e.GET("/htmx/sizes/:path", handlers.HandleSizes)

	e.Logger.Fatal(e.Start(":3000"))

}
