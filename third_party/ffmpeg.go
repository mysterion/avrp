package thirdparty

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/mysterion/avrp/internal/utils"
)

// returns [ffmpeg found in BUILD] [ffmpeg found in ENV]
func CheckFfmpeg() (bool, bool) {

	bins := []string{"ffmpeg", "ffprobe"}

	if runtime.GOOS == "windows" {
		for i, _ := range bins {
			bins[i] += ".exe"
		}
	}

	os.Stat(filepath.Join())

	checkPath, err := os.Executable()
	utils.Panic(err)

	if utils.DEV {
		checkPath, err = os.Getwd()
		utils.Panic(err)
	}

	buildC := 0
	envC := 0

	for i, _ := range bins {
		// fmt.Println(filepath.Join(checkPath, "third_party", bins[i]))
		if _, err = os.Stat(filepath.Join(checkPath, "third_party", bins[i])); err == nil {
			buildC += 1
		}

		if _, err := exec.LookPath(bins[i]); err == nil {
			envC += 1
		}

	}

	return buildC == len(bins), envC == len(bins)
}
