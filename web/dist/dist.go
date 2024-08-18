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

func Init() {
	VersionFile = filepath.Join(utils.ConfigDir, "VERSION")
}

// checks if dist is present
func Valid() bool {
	_, err := os.Stat(filepath.Join(utils.DistDir, "index.html"))
	return err == nil
}

func Delete() error {
	return os.RemoveAll(utils.DistDir)
}

// returns 0 , when no dist
func Ver() int {
	fd, err := os.Open(VersionFile)
	if errors.Is(err, fs.ErrNotExist) {
		log.Println("ERR: Version file doesn't exist")
		return 0
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
