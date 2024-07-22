package session

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"io"
	"nedas/shop/utils"
	"net/http"
)

// one problem if someone gets the session it will be valid for life
// later we will need to add some kind of created or when expires for more security
type Session struct {
	UserId string
	gcm    cipher.AEAD
	nonce  []byte
}

func (s *Session) String() string {
	hash := s.gcm.Seal(s.nonce, s.nonce, []byte(s.UserId), nil)
	return hex.EncodeToString(hash)
}

func (s *Session) UpdateNonce() {
	_, err := io.ReadFull(rand.Reader, s.nonce)
	if err != nil {
		panic(err)
	}
}

func (s *Session) Cookie() *http.Cookie {
	return &http.Cookie{
		Name:     "auth-session",
		Value:    s.String(),
		Path:     "/",
		MaxAge:   60 * 60 * 24 * 200,
		Secure:   false, // todo: how to know if we have ssh
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
}

func NewSession(userId string) *Session {
	secret := utils.Getenv("SESSION_SECRET")

	key, err := hex.DecodeString(secret)
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

	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		panic(err)
	}

	return &Session{
		UserId: userId,
		gcm:    gcm,
		nonce:  nonce,
	}
}

func SessionFromHash(sessionHash string) (*Session, bool) {
	if sessionHash == "" {
		return nil, false
	}

	hash, err := hex.DecodeString(sessionHash)
	if err != nil {
		return nil, false
	}

	secret := utils.Getenv("SESSION_SECRET")

	key, err := hex.DecodeString(secret)
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

	nonce, ciphertext := hash[:nonceSize], hash[nonceSize:]
	str, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, false
	}

	return &Session{
		UserId: string(str),
		gcm:    gcm,
		nonce:  nonce,
	}, true

}
