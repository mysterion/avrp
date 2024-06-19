package dist

import (
	"archive/zip"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/mysterion/avrp/internal/utils"
)

func DownloadLatest() error {

	r, err := latestRelease()
	if err != nil {
		return err
	}

	if len(r.Assets) == 0 {
		log.Fatal("Nothing to download in the latest release.")
	}

	var url string

	for _, a := range r.Assets {
		if strings.Contains(a.BrowserDownloadUrl, "dist") &&
			strings.HasSuffix(a.BrowserDownloadUrl, ".zip") {
			url = a.BrowserDownloadUrl
		}
	}

	if len(r.Assets) == 0 {
		log.Fatal("dist.zip not found in release.")
	}

	log.Printf("Latest Release - %v - %v\n", r.Tag, url)

	zipFile, err := os.CreateTemp("", "avrp-latest")
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

	err = extractZip(zipFile.Name(), utils.AppDir)
	if err != nil {
		return err
	}
	log.Println("Extracted successfully")

	log.Println("Updating `last updated` file")
	err = UpdateUFile()
	if err != nil {
		return err
	}

	return nil
}

func extractZip(srcZip string, dst string) error {
	log.Printf("Extracting %v to %v\n", srcZip, dst)
	r, err := zip.OpenReader(srcZip)
	if err != nil {
		return err
	}
	defer r.Close()

	err = os.MkdirAll(dst, 0755)
	if err != nil {
		return err
	}

	for _, f := range r.File {
		fi := f.FileInfo()
		target := filepath.Join(dst, f.Name)
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
