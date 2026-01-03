package thumbnails

import (
	"archive/zip"
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/mysterion/avrp/internal/utils"
)

var (
	ErrNotImpl    = errors.New("not implemented")
	ErrNoDownload = errors.New("no download found in the latest release")
)

var (
	binFfmpeg  string
	binFfprobe string
)

type release struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Tag    string `json:"tag_name"`
	URL    string `json:"html_url"`
	Assets []struct {
		ID                 int       `json:"id"`
		Size               int       `json:"size"`
		CreatedAt          time.Time `json:"created_at"`
		BrowserDownloadUrl string    `json:"browser_download_url"`
	} `json:"assets"`
}

func latestFfmpeg() (release, error) {
	repoOwner := "GyanD"
	repoName := "codexffmpeg"
	var r release

	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", repoOwner, repoName)

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return r, err
	}

	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return r, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return r, err
	}

	err = json.Unmarshal(body, &r)
	if err != nil {
		return r, err
	}

	return r, nil
}

func DownloadFfmpeg() error {
	if runtime.GOOS != "windows" || runtime.GOARCH != "amd64" {
		fmt.Println("\nPlease install ffmpeg from your package manager and make sure its in $PATH")
		return ErrNotImpl
	}

	r, err := latestFfmpeg()

	if err != nil {
		log.Printf("ERR: Fetching latest ffmpeg release - %s", err.Error())
		return err
	}

	if len(r.Assets) == 0 {
		return ErrNoDownload
	}

	url := ""

	for _, a := range r.Assets {
		if strings.Contains(a.BrowserDownloadUrl, "essentials") &&
			strings.HasSuffix(a.BrowserDownloadUrl, ".zip") {
			url = a.BrowserDownloadUrl
			break
		}
	}

	if url == "" {
		return ErrNoDownload
	}

	zipName := path.Base(url)
	zipPath := filepath.Join(os.TempDir(), zipName)

	zipFile, err := os.Create(zipPath)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	log.Println("Downloading to - ", zipPath)

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(zipFile, resp.Body)
	if err != nil {
		return err
	}

	err = extractFfmpeg(zipName, zipPath)
	if err != nil {
		log.Println("ERR: Failed to extract the zip")
		return err
	}

	log.Println("Extracted successfully")

	if !ThumbEnabled() {
		ThumbEnable()
	}

	return nil

}

func extractFfmpeg(zipName string, zipPath string) error {

	log.Printf("Extracting %v", zipPath)
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer r.Close()

	folderName := strings.TrimSuffix(zipName, filepath.Ext(zipName))

	for _, f := range r.File {
		fi := f.FileInfo()
		target := filepath.Join(utils.ConfigDir, strings.Replace(f.Name, folderName, "ffmpeg", 1))
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

func promptDownloadFfmpeg() bool {
	fmt.Println("\nffmpeg is required to show thumbnails in aframe-vr-player")
	fmt.Print("Download ffmpeg? \"yes\" or \"no\": ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	ans := strings.ToLower(strings.TrimSpace(scanner.Text()))
	return ans == "yes"
}

func initFfmpeg() bool {

	ffmpeg := "ffmpeg"
	ffprobe := "ffprobe"

	if runtime.GOOS == "windows" {
		ffmpeg += ".exe"
		ffprobe += ".exe"
	}

	var err1, err2 error

	binFfmpeg, err1 = exec.LookPath(ffmpeg)
	binFfprobe, err2 = exec.LookPath(ffprobe)

	if err1 == nil && err2 == nil {
		log.Println("ffmpeg found in $PATH")
		return true
	}

	binFfmpeg = filepath.Join(utils.FfmpegDir, "bin", ffmpeg)
	binFfprobe = filepath.Join(utils.FfmpegDir, "bin", ffprobe)

	_, err1 = os.Stat(binFfmpeg)
	_, err2 = os.Stat(binFfprobe)

	if err1 == nil && err2 == nil {
		log.Println("ffmpeg found in the ffmpegDir")
		return true
	}

	// code to download ffmpeg download 👇👇👇👇
	accept := promptDownloadFfmpeg()

	if !accept {
		fmt.Printf("\n\nYou can disable this message, by running: avrp --no-thumb\n\n")
		return false
	}

	err := DownloadFfmpeg()

	if err != nil {
		log.Println("ERR: Failed to download ffmpeg")
		return false
	}

	return true
}

func ThumbEnable() {
	err := os.Remove(utils.ThumbDisabledFile)
	utils.Panic(err)
}

func ThumbDisable() {
	_, err := os.Create(utils.ThumbDisabledFile)
	utils.Panic(err)
}

func ThumbEnabled() bool {
	_, err := os.Stat(utils.ThumbDisabledFile)
	return err != nil
}
