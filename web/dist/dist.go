package dist

import (
	"errors"
	"io/fs"
	"log"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/mysterion/avrp/internal/utils"
)

var VersionFile string

func init() {
	VersionFile = filepath.Join(utils.DistDir, "version")
}

// checks if dist is present
func Valid() bool {
	_, err := os.Stat(filepath.Join(utils.DistDir, "index.html"))
	return err == nil
}

func Delete() error {
	return os.Remove(utils.DistDir)
}

// returns 0 , when no dist
//
// returns math.MaxInt, when no Version file found(maybe custom version?)
func Ver() int {
	fd, err := os.Open(VersionFile)
	if errors.Is(err, fs.ErrNotExist) {

		log.Println("ERR: Version file doesn't exist")

		if !Valid() {
			return 0
		}

		return math.MaxInt
	}
	defer fd.Close()

	d := make([]byte, 1024)
	n, err := fd.Read(d)
	if err != nil {
		return math.MaxInt
	}

	vs := strings.ReplaceAll(string(d[:n]), ".", "")
	v, err := strconv.Atoi(vs)
	if err != nil {
		return math.MaxInt
	}

	return v
}
