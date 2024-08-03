package utils

import "os"

func Getenv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		panic(key + " is not set")
	}
	return value
}

func Assert(ok bool, msg string) {
	if !ok {
		panic(msg)
	}
}

func TLSEnabled() bool {
	// we will do it by reading env vars and checking if there is certifiates
	return false
}
