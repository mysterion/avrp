package dist

import (
	"log"
	"os"
	"path/filepath"

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

// returns sha of aframe-vr-player dist
func Ver() string {
	d, err := os.ReadFile(VersionFile)

	if err != nil {
		log.Printf("WARN: while reading dist version, %v\n", err.Error())
		return ""
	}

	return string(d)
}
