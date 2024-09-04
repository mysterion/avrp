package dist

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/mysterion/avrp/internal/utils"
)

func DownloadCommit(sha string) error {

	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/zipball/%s", RepoOwner, RepoName, sha)
	log.Printf("Downloading - %v\n", url)

	zipFile, err := os.CreateTemp("", "aframe-vr-player")
	if err != nil {
		return err
	}
	defer zipFile.Close()

	log.Println("Downloading to - ", zipFile.Name())

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(zipFile, resp.Body)
	if err != nil {
		return err
	}

	err = extractZip(sha, zipFile.Name())
	if err != nil {
		log.Println("ERR: Failed to extract the zip")
		return err
	}

	err = os.WriteFile(VersionFile, []byte(sha), 0644)

	if err != nil {
		log.Println("ERR: Failed to write version file")
		return err
	}

	log.Println("Extracted successfully")

	log.Println("Updating `last updated` file")

	return nil
}

func extractZip(sha string, srcZip string) error {

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

	folderName := fmt.Sprintf("%s-%s-%s", RepoOwner, RepoName, sha[:7])

	for _, f := range r.File {
		fi := f.FileInfo()
		target := filepath.Join(utils.ConfigDir, strings.Replace(f.Name, folderName, "dist", 1))
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

	return nil
}
