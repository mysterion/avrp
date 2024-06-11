package utils

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

var ConfigDir string
var AppDir string

var DEV bool

func init() {
	DEV = len(os.Getenv("DEV")) > 0

	h, err := os.UserHomeDir()
	Panic(err)

	ConfigDir = filepath.Join(h, ".avrp")

	AppDir, err = os.Executable()
	Panic(err)
	AppDir = filepath.Dir(AppDir)
	if DEV {
		log.Println("***RUNNING AS DEV***")
		AppDir, err = os.Getwd()
		Panic(err)
	}

	err = os.MkdirAll(ConfigDir, 0755)
	Panic(err)
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
