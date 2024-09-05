package utils

import (
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

var ConfigDir string
var DistDir string
var AppDir string

var DEV bool

func Init() {

	h, err := os.UserHomeDir()
	Panic(err)

	ConfigDir = filepath.Join(h, ".avrp")

	DistDir = filepath.Join(ConfigDir, "dist")
	if DEV {
		DistDir = "."
	}

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

	err = os.MkdirAll(DistDir, 0755)
	Panic(err)

	if DEV {
		log.Printf("Home directory: %s\n", h)
		log.Printf("Config directory: %s\n", ConfigDir)
		log.Printf("Dist directory: %s\n", DistDir)
		log.Printf("App Run directory: %s\n", AppDir)
	}

}

func Panic(err error) {
	_, filename, lno, _ := runtime.Caller(1)
	if err != nil {
		log.Printf("Panic call from : %v:%v\n", filename, lno)
		panic(err)
	}
}

func IsGoRun() bool {
	return strings.HasPrefix(AppDir, filepath.Join(os.TempDir(), "go-build"))
}
