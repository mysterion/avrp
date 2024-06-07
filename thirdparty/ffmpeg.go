package thirdparty

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/mysterion/avrp/internal/utils"
)

// empty if ffmpeg is not found in system
var FfmpegBin = ""

// empty if ffprobe is not found in system
var FfprobeBin = ""

func init() {

	ffmpeg := "ffmpeg"
	ffprobe := "ffprobe"

	if runtime.GOOS == "windows" {
		ffmpeg += ".exe"
		ffprobe += ".exe"
	}

	checkPath, err := os.Executable()
	utils.Panic(err)

	if utils.DEV {
		checkPath, err = os.Getwd()
		utils.Panic(err)
	}

	// check in distribution
	FfmpegBin = filepath.Join(checkPath, "thirdparty", ffmpeg)
	_, err = os.Stat(FfmpegBin)
	if err != nil {
		log.Println(err)
		FfmpegBin = ""
	}
	FfprobeBin = filepath.Join(checkPath, "thirdparty", ffprobe)
	_, err = os.Stat(FfprobeBin)
	if err != nil {
		log.Println(err)
		FfprobeBin = ""
	}

	// check in PATH
	if FfmpegBin == "" {
		ffmpegPath, err := exec.LookPath(ffmpeg)
		log.Println(err)
		if err == nil {
			FfmpegBin = ffmpegPath
		} else {
			log.Println(err)
		}

	}
	if FfprobeBin == "" {
		ffprobePath, err := exec.LookPath(ffprobe)

		if err == nil {
			FfprobeBin = ffprobePath
		} else {
			log.Println(err)
		}
	}

}
