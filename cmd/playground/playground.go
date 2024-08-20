package main

import (
	"log"
)

func main() {
	log.SetFlags(log.Default().Flags() | log.Lshortfile)
}
