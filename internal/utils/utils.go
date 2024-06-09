package utils

import (
	"os"
	"path/filepath"
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
	if DEV {
		AppDir, err = os.Getwd()
		Panic(err)
	}

	err = os.MkdirAll(ConfigDir, os.ModeDir)
	Panic(err)
}

func Panic(err error) {
	if err != nil {
		panic(err)
	}
}
