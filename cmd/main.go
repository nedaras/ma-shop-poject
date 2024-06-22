package main

import (
  "github.com/labstack/echo/v4"

	"nedas/shop/pkg/handlers"
)

func main() {

	e := echo.New()

  e.Static("/", "public");

	e.GET("/", handlers.HandleIndex)

	e.Logger.Fatal(e.Start(":3000"))

}
