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
