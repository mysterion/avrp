package main

import (
	"fmt"
	"log"

	thirdparty "github.com/mysterion/avrp/third_party"
)

func main() {
	// api.Start(5000)
	log.SetFlags(log.Lshortfile)
	fmt.Println(thirdparty.CheckFfmpeg())
}
