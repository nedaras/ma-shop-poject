package main

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"nedas/shop/pkg/handlers"
	"nedas/shop/pkg/middlewares"
	"nedas/shop/pkg/models"
	"nedas/shop/pkg/storage"
)

func main() {

	c, err := storage.NewCassandra()
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	user := models.User{
		UserID: "123456",
		Email:  "pimpalas.gaidys@gmail.com",
	}

	u2, err := c.GetUser(user.UserID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(u2)

	return

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
