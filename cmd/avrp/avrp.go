package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/mysterion/avrp/web/api"
)

func listFilesAndFolders(dirPath string) (error, []string, []string) {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return err, nil, nil
	}
	files := make([]string, 0)
	folders := make([]string, 0)
	for _, entry := range entries {
		if entry.IsDir() {
			folders = append(folders, entry.Name())
		} else {
			files = append(files, entry.Name())
		}
	}
	return nil, files, folders
}

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
