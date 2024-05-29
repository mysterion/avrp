package main

import (
	"fmt"

	"github.com/mysterion/avrp/web"
)

func main() {
	fmt.Println("Hello world")
	r, err := web.LatestRelease()
	if err != nil {
		panic(err)
	}
	fmt.Println("Latest Release: ", r.Tag)

	fmt.Println("All Releases: ")

	rall, err := web.AllReleases()
	if err != nil {
		panic(err)
	}

	for _, release := range rall {
		fmt.Println(release.Tag)
	}

}
