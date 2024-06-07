package thumbnails

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/mysterion/avrp/internal/utils"
	"github.com/mysterion/avrp/thirdparty"
)

var thumbdir string

var Available = true

func init() {
	thumbdir = filepath.Join(utils.ConfigDir, "thumbnails")
	err := os.MkdirAll(thumbdir, os.ModeDir)

	if err != nil {
		log.Printf("Thumbnails not available - %v", err)
		Available = false
	}

	if thirdparty.FfmpegBin == "" || thirdparty.FfprobeBin == "" {
		log.Printf("Thumbnails not available - ffmpeg not found")
		Available = false
	}
}

func GetDuration(file string) (float64, error) {
	var seconds float64
	args := strings.Split(fmt.Sprintf("-v error -show_entries format=duration -of default=noprint_wrappers=1:nokey=1 %s", file), " ")
	cmd := exec.Command(thirdparty.FfprobeBin, args...)
	stdout, err := cmd.Output()
	if err != nil {
		log.Println("Failed to get Duration - ", err)
		return seconds, err
	}
	return strconv.ParseFloat(strings.TrimSpace(string(stdout)), 64)
}

func Generate(file string) {

}
