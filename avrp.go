package main

import (
	"bufio"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"

	"github.com/mysterion/avrp/internal/utils"
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

var servDir string

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	utils.GoRunGatekeeper()

	if !utils.DEV {
		dist.TryUpdate()
	}

	args := os.Args[1:]
	if len(args) == 0 {
		servDir = askForPath()
	} else {
		if isPathValid(args[0]) {
			servDir = args[0]
		} else {
			servDir = askForPath()
		}
	}
	api.Init(servDir)
	api.Start(5000)
}
