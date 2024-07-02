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

// checks if dist is present
func Valid() bool {
	if utils.DEV {
		return true
	}
	_, err := os.Stat(filepath.Join(utils.AppDir, "dist", "index.html"))
	return err == nil
}

func Delete() error {
	return os.Remove(filepath.Join(utils.AppDir, "dist"))
}

// returns 0 , when no dist
//
// returns math.MaxInt, when no version file can be found on dist
//
// indicating a non-standard version
func Ver() int {
	fd, err := os.Open(filepath.Join(utils.AppDir, "dist", "version"))
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
