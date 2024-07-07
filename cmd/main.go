package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"nedas/shop/pkg/handlers"
	"nedas/shop/pkg/middlewares"
)

type Session struct {
  UserId string
  hash []byte
}

func (s *Session) String() string {
  return hex.EncodeToString(s.hash)
}

func newSession(userId string) *Session {
  key, err := hex.DecodeString(os.Getenv("SESSION_SECRET"))
  if err != nil {
    panic(err)
  }

  block, err := aes.NewCipher(key)
  if err != nil {
    panic(err)
  }

  gcm, err := cipher.NewGCM(block)
  if err != nil {
    panic(err)
  }
  
  nonce := make([]byte, gcm.NonceSize())
  io.ReadFull(rand.Reader, nonce)

  // bad what us we change sessions id and then call string boom, we explode
  return &Session{
    UserId: userId,
    hash: gcm.Seal(nonce, nonce, []byte(userId), nil),
  }
}

func sessionFromHash(hash string) (*Session, bool) {
  bh, err := hex.DecodeString(hash)
  if err != nil {
    return nil, false
  }

  key, err := hex.DecodeString(os.Getenv("SESSION_SECRET"))
  if err != nil {
    panic(err)
  }

  block, err := aes.NewCipher(key)
  if err != nil {
    panic(err)
  }

  gcm, err := cipher.NewGCM(block)
  if err != nil {
    panic(err)
  }

  nonceSize := gcm.NonceSize()
  if nonceSize > len(hash) {
    return nil, false
  }

  nonce, ciphertext := bh[:nonceSize], bh[nonceSize:]
	str, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, false 
	}

  return &Session{
    UserId: string(str),
    hash: bh,
  }, true

}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

  session := newSession("some data that is static for now")
  fmt.Println(session.String())

  session, ok := sessionFromHash(session.String())
  if ok {
    fmt.Println(session.UserId)
  }

  return

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

	e.Logger.Fatal(e.Start(":3000"))
}
