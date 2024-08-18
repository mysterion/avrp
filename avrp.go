package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"

	"github.com/mysterion/avrp/internal/thumbnails"
	"github.com/mysterion/avrp/internal/utils"
	"github.com/mysterion/avrp/thirdparty"
	"github.com/mysterion/avrp/web/api"
	"github.com/mysterion/avrp/web/dist"
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
	var update bool
	var reset bool

	flag.BoolVar(&utils.DEV, "dev", false, "starts in dev mode, serves 'index.html' from current directory")
	flag.BoolVar(&update, "update", false, "checks & downloads the latest version of 'aframe-vr-player'")
	flag.StringVar(&servDir, "dir", "", "path to video files")
	flag.BoolVar(&reset, "reset", false, "removes all configs, thumbnails & 'aframe-vr-player' files")

	flag.Parse()

	utils.Init()
	dist.Init()
	thirdparty.Init()
	thumbnails.Init()

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
		dist.Update()
	}

	/// do i need this anymore?
	utils.GoRunGatekeeper()

	// running for the first time
	if !dist.Valid() && !utils.DEV {
		dist.Update()
	}

	if len(servDir) == 0 {
		servDir = askForPath()
	} else {
		if !isPathValid(servDir) {
			servDir = askForPath()
		}
	}

	api.Init(servDir)
	api.Start(5000)
}
