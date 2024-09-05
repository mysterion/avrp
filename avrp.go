package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"

	"github.com/mysterion/avrp/internal/dist"
	"github.com/mysterion/avrp/internal/server"
	"github.com/mysterion/avrp/internal/thumbnails"
	"github.com/mysterion/avrp/internal/utils"
)

func isPathValid(path string) bool {
	_, err := os.Stat(path)
	return !errors.Is(err, fs.ErrNotExist)
}

func getUserInput(prompt string) string {
	fmt.Println(prompt + ": ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return scanner.Text()
}

func askForPath() string {
	for {
		path := getUserInput("Please enter path to serve")
		if isPathValid(path) {
			return path
		} else {
			fmt.Println("Invalid Path")
		}
	}
}

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	var servDir string
	var sha string
	var update bool
	var reset bool
	var port int
	var getFfmpeg bool
	var noThumb bool

	flag.BoolVar(&utils.DEV, "dev", false, "starts in dev mode, serves 'index.html' from current directory")
	flag.BoolVar(&update, "update", false, "checks & downloads the latest version(commit) of 'aframe-vr-player'")
	flag.StringVar(&sha, "sha", "latest", "Optional - download a specific commit of aframe-vr-player")
	flag.StringVar(&servDir, "dir", "", "path to video files")
	flag.BoolVar(&reset, "reset", false, "removes all configs, thumbnails & 'aframe-vr-player' files")
	flag.IntVar(&port, "port", 5000, "port to serve on (default 5000)")
	flag.BoolVar(&getFfmpeg, "get-ffmpeg", false, "downloads ffmpeg")
	flag.BoolVar(&noThumb, "no-thumb", false, "disables thumbnail generation")

	flag.Parse()

	utils.Init()
	dist.Init()
	thumbnails.Init()

	// commands 👇

	if getFfmpeg {
		utils.Panic(thumbnails.DownloadFfmpeg())
		return
	}

	if noThumb {
		if thumbnails.NoFfmpeg() {
			thumbnails.NoFfmpegFileRemove()
			log.Println("thumbnail generation enabled👍")
		} else {
			thumbnails.NoFfmpegFileCreate()
			log.Println("thumbnail generation disabled👎")
		}
		return
	}

	if reset {
		err := os.RemoveAll(utils.ConfigDir)
		if err != nil {
			log.Panicln(err)
		} else {
			log.Printf("Success! removed dir : %s\n", utils.ConfigDir)
		}
		return
	}

	if update {
		dist.Update(sha)
		return
	}

	// commands 👆

	// running for the first time
	if !dist.Valid() && !utils.DEV {
		dist.Update("latest")
	}

	if len(servDir) == 0 {
		servDir = askForPath()
	} else {
		if !isPathValid(servDir) {
			servDir = askForPath()
		}
	}

	server.Init(servDir)
	server.Start(port)
}
