package thirdparty

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/mysterion/avrp/internal/logg"
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

	// check in distribution
	FfmpegBin = filepath.Join(utils.AppDir, "thirdparty", ffmpeg)
	_, err := os.Stat(FfmpegBin)
	if err != nil {
		logg.Debug(err)
		FfmpegBin = ""
	}
	FfprobeBin = filepath.Join(utils.AppDir, "thirdparty", ffprobe)
	_, err = os.Stat(FfprobeBin)
	if err != nil {
		logg.Debug(err)
		FfprobeBin = ""
	}

	// check in PATH
	if FfmpegBin == "" {
		ffmpegPath, err := exec.LookPath(ffmpeg)
		if err == nil {
			FfmpegBin = ffmpegPath
		} else {
			logg.Debug(err)
		}

	}
	if FfprobeBin == "" {
		ffprobePath, err := exec.LookPath(ffprobe)

		if err == nil {
			FfprobeBin = ffprobePath
		} else {
			logg.Debug(err)
		}
	}

	logg.Debug("Ffmpeg", FfmpegBin)
	logg.Debug("Ffprobe", FfprobeBin)
}
