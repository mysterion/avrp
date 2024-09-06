package utils

import (
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

var (
	ConfigDir string
	DistDir   string
	AppDir    string
	FfmpegDir string
	ThumbDir  string
)

var (
	NoThumbFile string
	VersionFile string
)

var DEV bool

func Init() {

	h, err := os.UserHomeDir()
	Panic(err)

	ConfigDir = filepath.Join(h, ".avrp")

	FfmpegDir = filepath.Join(ConfigDir, "ffmpeg")

	ThumbDir = filepath.Join(ConfigDir, "thumbnails")

	NoThumbFile = filepath.Join(ConfigDir, "nothumb")

	VersionFile = filepath.Join(ConfigDir, "VERSION")

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

	Panic(os.MkdirAll(ConfigDir, 0755))

	Panic(os.MkdirAll(DistDir, 0755))

	Panic(os.MkdirAll(FfmpegDir, 0755))

	Panic(os.MkdirAll(ThumbDir, 0755))

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
