package thumbnails

import (
	"errors"
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
)

var thumbdir string

var ErrNotVideo = errors.New("not a video")

var Available = false

var muGen sync.Mutex

func Init() {
	noffmpegfile = filepath.Join(utils.ConfigDir, "noffmpeg")

	thumbdir = filepath.Join(utils.ConfigDir, "thumbnails")
	err := os.MkdirAll(thumbdir, 0755)
	if err != nil {
		log.Printf("ERR: Thumbnails not available - %v\n", err)
		return
	}

	ffmpegDir = filepath.Join(utils.ConfigDir, "ffmpeg")
	err = os.MkdirAll(ffmpegDir, 0755)
	if err != nil {
		log.Printf("ERR: Thumbnails not available - %v\n", err)
		return
	}

	if NoFfmpeg() {
		return
	}

	var found = false
	found, binFfmpeg, binFfprobe = CheckFfmpegInPath()
	if found {
		Available = true
		return
	}

	found, binFfmpeg, binFfprobe = CheckFfmpeg()
	if found {
		Available = true
		return
	}

	accept := promptDownloadFfmpeg()

	if !accept {
		fmt.Printf("\n\nYou can disable this message, by running: avrp --no-thumb\n\n")
		return
	}

	utils.Panic(DownloadFfmpeg())

	found, binFfmpeg, binFfprobe = CheckFfmpeg()
	if found {
		Available = true
		return
	} else {
		log.Println("Something went wrong, please re-download ffmpeg: avrp --get-ffmpeg")
	}
}

func GetDuration(file string) (float64, error) {
	if !utils.IsVideo(file) {
		return 0, ErrNotVideo
	}
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
		cmd := exec.Command(binFfprobe, args...)
		stdout, err := cmd.CombinedOutput()
		if err != nil {
			log.Printf("ERR - Failed to get Duration for %v - %v\nSTDOUT:\n%s\n", file, err, stdout)
			stdout = []byte("0")
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
	err = os.MkdirAll(outDir, 0755)

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
			cmd := exec.Command(binFfmpeg, cmdArgs...)
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
