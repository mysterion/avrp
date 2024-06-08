package thumbnails

import (
	"fmt"
	"log"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/mysterion/avrp/internal/cache"
	"github.com/mysterion/avrp/internal/utils"
	"github.com/mysterion/avrp/thirdparty"
)

var thumbdir string

var Available = true

var muGen sync.Mutex

func init() {
	thumbdir = filepath.Join(utils.ConfigDir, "thumbnails")
	err := os.MkdirAll(thumbdir, os.ModeDir)

	if err != nil {
		log.Printf("Thumbnails not available - %v\n", err)
		Available = false
	}

	if thirdparty.FfmpegBin == "" || thirdparty.FfprobeBin == "" {
		log.Printf("Thumbnails not available - ffmpeg not found\n")
		Available = false
	}
}

func GetDuration(file string) (float64, error) {
	var secs string
	secs = cache.Get("DUR_" + file)
	if secs == "" {
		args := []string{
			"-v", "error",
			"-show_entries",
			"format=duration",
			"-of", "default=noprint_wrappers=1:nokey=1",
			file,
		}
		cmd := exec.Command(thirdparty.FfprobeBin, args...)
		stdout, err := cmd.CombinedOutput()
		if err != nil {
			log.Printf("ERR - Failed to get Duration - %v\nSTDOUT:\n%s\n", err, stdout)
			return 0, err
		}
		secs = strings.TrimSpace(string(stdout))

		cache.Set("DUR_"+file, secs)
	}

	return strconv.ParseFloat((secs), 64)
}

func Generated(file string) bool {

	h, err := Hash(file)

	if err != nil {
		log.Println("ERR - ", err)
		return false
	}

	duration, err := GetDuration(file)
	if err != nil {
		log.Println("ERR - ", err)
		return false
	}

	p := filepath.Join(thumbdir, h, fmt.Sprintf("%v.jpg", math.Floor(duration/60)-1))
	_, err = os.Stat(p)

	return err == nil
}

// TODO: keep error state for a particular file with eviction policy
func Generate(file string) {
	muGen.Lock()
	defer muGen.Unlock()
	if Generated(file) {
		log.Printf("Already Generated - %v\n", file)
		return
	}
	h, err := Hash(file)
	if err != nil {
		log.Printf("ERR: %v\n", err.Error())
	}

	outDir := filepath.Join(thumbdir, h)
	err = os.MkdirAll(outDir, os.ModeDir)

	if err != nil {
		log.Printf("ERR while creating thumbnail dir: %v\n", err.Error())
		return
	}

	duration, err := GetDuration(file)
	if err != nil {
		log.Printf("ERR while getting duration of the input file: %v\n", err.Error())
		return
	}
	n := int(math.Floor(duration / 60))
	done := make(chan bool, n)
	defer close(done)

	for i := 0; i < n; i++ {
		go func(i int, done chan<- bool) {
			defer func() { done <- true }()
			cmdArgs := []string{
				"-y", "-accurate_seek", "-ss", fmt.Sprintf("%v", i*60),
				"-i", file,
				"-frames:v", "1",
				"-vf", "crop=in_w/2:in_h/2:in_w:in_h/4,scale=320:-1",
				filepath.Join(outDir, fmt.Sprintf("%v.jpg", i)),
			}
			cmd := exec.Command(thirdparty.FfmpegBin, cmdArgs...)
			log.Println("EXEC - ", thirdparty.FfmpegBin, cmdArgs)
			stdout, err := cmd.CombinedOutput()
			if err != nil {
				log.Printf("ERR while generating thumbnail %vth for %v - %v\nSTDOUT:\n%v\n", i, file, err, string(stdout))
				return
			}

		}(i, done)

	}
	for i := 0; i < n; i++ {
		<-done
	}
}

func Get(id string, file string) (string, error) {
	h, err := Hash(file)
	if err != nil {
		return "", err
	}
	return filepath.Join(thumbdir, h, fmt.Sprintf("%v.jpg", id)), nil
}
