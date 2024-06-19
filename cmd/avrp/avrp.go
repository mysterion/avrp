package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/mysterion/avrp/internal/utils"
	"github.com/mysterion/avrp/web/api"
	"github.com/mysterion/avrp/web/dist"
)

func isPathValid(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
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
	utils.GoRunGatekeeper()

	update, err := dist.CheckUFile()
	if err != nil {
		log.Println("ERR: while checking for update, ", err)
		log.Println("Skipping Update check")
	}

	if update || !dist.Ok() {
		log.Println("Downloading Latest version of Aframe Vr Player")
		err := dist.DownloadLatest()
		utils.Panic(err)
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
