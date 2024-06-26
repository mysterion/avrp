package dist

import (
	"os"
	"path/filepath"

	"github.com/mysterion/avrp/internal/utils"
)

func Ok() bool {
	if utils.DEV {
		return true
	}
	_, err := os.Stat(filepath.Join(utils.AppDir, "dist", "index.html"))
	return err == nil
}
