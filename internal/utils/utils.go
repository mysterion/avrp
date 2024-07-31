package utils

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

var ConfigDir string
var DistDir string
var AppDir string
var UpdateFile string

var DEV bool

func init() {
	DEV = len(os.Getenv("DEV")) > 0

	h, err := os.UserHomeDir()
	Panic(err)

	ConfigDir = filepath.Join(h, ".avrp")

	DistDir = filepath.Join(ConfigDir, "dist")

	AppDir, err = os.Executable()
	Panic(err)

	AppDir = filepath.Dir(AppDir)
	if DEV {
		log.Println("***RUNNING AS DEV***")
		AppDir, err = os.Getwd()
		Panic(err)
	}

	UpdateFile = filepath.Join(ConfigDir, "LAST_UPDATE_CHECK")

	err = os.MkdirAll(ConfigDir, 0755)
	Panic(err)

	err = os.MkdirAll(DistDir, 0755)
	Panic(err)

	if DEV {
		log.Printf("Home directory: %s\n", h)
		log.Printf("Config directory: %s\n", ConfigDir)
		log.Printf("Dist directory: %s\n", DistDir)
		log.Printf("App Run directory: %s\n", AppDir)
		log.Printf("Update File: %s\n", UpdateFile)
	}

}

func Panic(err error) {
	if err != nil {
		panic(err)
	}
}

func GoRunGatekeeper() {
	if strings.HasPrefix(AppDir, filepath.Join(os.TempDir(), "go-build")) {
		panic("MAYBE YOU FORGOT DEV=1 ? I'M NOT LETTING YOU RUN STUFF FROM TEMP DIRECTORY")
	}

}
