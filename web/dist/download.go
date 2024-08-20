package dist

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/mysterion/avrp/internal/utils"
)

func DownloadTag(t Tag) error {

	log.Printf("Latest Release - %v - %v\n", t.Name, t.ZipballUrl)

	zipFile, err := os.CreateTemp("", "avrp-latest")
	if err != nil {
		return err
	}
	defer zipFile.Close()

	log.Println("Downloading to - ", zipFile.Name())

	resp, err := http.Get(t.ZipballUrl)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(zipFile, resp.Body)
	if err != nil {
		return err
	}

	err = extractZip(t, zipFile.Name())
	if err != nil {
		log.Println("ERR: Failed to extract the zip")
		return err
	}

	err = os.WriteFile(VersionFile, []byte(t.Name), 0644)

	if err != nil {
		log.Println("ERR: Failed to write version file")
		return err
	}

	log.Println("Extracted successfully")

	log.Println("Updating `last updated` file")

	return nil
}

func extractZip(t Tag, srcZip string) error {
	// extracts ~/.avrp/<release-folder-sha>
	// renames ~/.avrp/<release-folder-sha> to ~/.avrp/dist
	log.Printf("Extracting %v", srcZip)
	r, err := zip.OpenReader(srcZip)
	if err != nil {
		return err
	}
	defer r.Close()

	err = os.MkdirAll(utils.ConfigDir, 0755)
	if err != nil {
		return err
	}

	for _, f := range r.File {
		fi := f.FileInfo()
		target := filepath.Join(utils.ConfigDir, f.Name)
		log.Println("Extracting - ", target)
		if !fi.IsDir() {
			err := os.MkdirAll(filepath.Dir(target), 0755)
			if err != nil {
				return err
			}

			file, err := os.Create(target)
			if err != nil {
				return err
			}

			eFile, err := f.Open()
			if err != nil {
				return err
			}
			defer eFile.Close()

			_, err = io.Copy(file, eFile)
			if err != nil {
				return err
			}
		}
	}

	if !utils.DEV {
		err = os.RemoveAll(utils.DistDir)
		if err != nil {
			return err
		}
	}

	newDist := filepath.Join(utils.ConfigDir, fmt.Sprintf("%s-%s-%s", RepoOwner, RepoName, t.Commit.Sha[:7]))

	err = os.Rename(newDist, utils.DistDir)
	if err != nil {
		return err
	}

	return nil
}
