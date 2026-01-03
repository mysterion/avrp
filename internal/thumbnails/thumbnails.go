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

	"github.com/mysterion/avrp/internal/cache"
	"github.com/mysterion/avrp/internal/utils"
)

var ErrNotVideo = errors.New("not a video")

var Available = false

var (
	guardFile  GuardFile
	guardProcs GuardProcs
)

func Init() {

	Available = initFfmpeg()

	if !Available {
		return
	}

	guardFile = NewGuardFile()
	guardProcs = NewGuardProcs()

}

func GetDuration(file string) (float64, error) {
	if !Available || !utils.IsVideo(file) {
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
	guardFile.Lock(file)
	defer guardFile.Unlock(file)

	if cache.Get("GEN_"+file) != "" {
		return true
	}

	if !Available || !utils.IsVideo(file) {
		return false
	}

	h, err := Hash(file)

	if err != nil {
		return false
	}

	duration, err := GetDuration(file)
	if err != nil {
		return false
	}

	count := math.Floor(duration / 60)

	for i := 0; i < int(count); i++ {
		t := fmt.Sprintf("%d.jpg", i)
		fd, err := os.Stat(filepath.Join(utils.ThumbDir, h, t))
		if err != nil || fd.Size() == 0 {
			return false
		}
	}

	cache.Set("GEN_"+file, "OK")
	return true
}

func Generate(file string) {
	guardFile.Lock(file)
	defer guardFile.Unlock(file)

	h, err := Hash(file)
	if err != nil {
		log.Printf("ERR: while generating file hash - %v\n", err.Error())
	}

	outDir := filepath.Join(utils.ThumbDir, h)
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
		guardProcs.In()
		go func(i int, done chan<- bool) {
			defer guardProcs.Out()
			defer func() { done <- true }()
			cmdArgs := []string{
				"-y", "-accurate_seek", "-ss", fmt.Sprintf("%v", i*60),
				"-i", file,
				"-frames:v", "1",
				"-vf", "crop=in_w/2:in_h/2:in_w:in_h/4,scale=320:-1",
				filepath.Join(outDir, fmt.Sprintf("%v.jpg", i)),
			}
			cmd := exec.Command(binFfmpeg, cmdArgs...)
			log.Printf("generating thumbnail %v.jpg :::: %v\n", i, filepath.Base(file))
			stdout, err := cmd.CombinedOutput()
			if err != nil {
				log.Printf("ERR while generating thumbnail %v.jpg :::: %v - %v\nSTDOUT:\n%v\n", i, file, err, string(stdout))
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
	return filepath.Join(utils.ThumbDir, h, fmt.Sprintf("%v.jpg", id)), nil
}

