package utils

import (
	"os"
)

var DEV bool

func init() {
	DEV = len(os.Getenv("DEV")) > 0
}

func Panic(err error) {
	if err != nil {
		panic(err)
	}
}
