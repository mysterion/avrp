package main

import (
	"github.com/mysterion/avrp/internal/thumbnails"
)

func main() {
	// api.Start(5000)
	// log.SetFlags(log.Lshortfile)
	// fmt.Println(thirdparty.FfmpegBin)
	// fmt.Println(thirdparty.FfprobeBin)
	// fmt.Println(thumbnails.Available)
	// fmt.Println(thumbnails.GetDuration("c:\\Users\\User\\Downloads\\Videos\\sample.mp4"))
	thumbnails.Generate("c:\\Users\\User\\Downloads\\Videos\\hm-dis.mp4")
}
